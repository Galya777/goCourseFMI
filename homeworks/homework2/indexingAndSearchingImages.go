package homework2

// go-web-crawler-image-indexer
// Single-file example.go implementation of the requested program.
// Features implemented:
//  - Recursive crawler starting from one or more URLs provided as command-line args
//  - Worker pool (configurable size) implemented with goroutines and channels
//  - Option to follow external links (flag -follow-external)
//  - Crawling timeout (flag -timeout, default 2m)
//  - Max concurrent goroutines limit (flag -max-goroutines)
//  - Headless browser support via chromedp to render JS single-page apps (flag -enable-js)
//  - Image extraction (raster formats and SVG). Raster thumbnails are generated (max width 200px).
//  - SVG files are saved; rasterizing SVG to PNG thumbnails is optional via external tool (see notes).
//  - Image metadata stored in MySQL (configurable via DSN flag)
//  - Small HTTP server with HTML templates for searching and visualizing images
//
// Limitations / Notes:
//  - SVG rasterization is not implemented in pure go here: to generate PNG thumbnails from SVG
//    you may provide an external rasterizer command (e.g. `rsvg-convert` or `inkscape`) via flag
//    -svg-rasterize-cmd. If omitted, SVG thumbnails won't be raster-generated; the SVG file will
//    still be saved and indexed.
//  - This is an example.go / reference implementation and includes minimal error handling for
//    clarity. For production use, add robust error handling, logging, retries, politeness
//    (robots.txt, rate-limiting), TLS verification options, etc.
//
// Dependencies (go get):
//  go get github.com/chromedp/chromedp
//  go get github.com/go-sql-driver/mysql
//  go get golang.org/x/net/html
//  go get golang.org/x/net/publicsuffix
//
// Build:
//  go build -o crawler main.go
//
// Example run:
//  ./crawler -workers=10 -timeout=2m -follow-external=false -enable-js=true -image-dir=images \
//      -mysql-dsn="user:pass@tcp(localhost:3306)/imagedb?parseTime=true" https://example.com
//
// Database schema (MySQL):
//
// CREATE DATABASE imagedb CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
// USE imagedb;
// CREATE TABLE images (
//   id BIGINT AUTO_INCREMENT PRIMARY KEY,
//   url TEXT NOT NULL,
//   filename VARCHAR(1024) NOT NULL,
//   thumbnail_path VARCHAR(1024),
//   alt_text VARCHAR(1024),
//   title_text VARCHAR(1024),
//   width INT,
//   height INT,
//   format VARCHAR(50),
//   crawled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
//
// High-level design notes:
//  - A dispatcher goroutine accepts starting URLs and keeps a "to visit" queue.
//  - Worker goroutines fetch pages (optionally using chromedp for JS rendering), parse links
//    and images, and send discovered links back to dispatcher to be scheduled if not seen.
//  - A visited map with sync.Mutex prevents revisiting.
//  - A semaphore (buffered channel) limits maximum concurrent HTTP fetch goroutines.
//  - Image downloads are performed by workers and thumbnails are created using the image
//    packages (jpeg/png/gif). Metadata is inserted into MySQL.
//

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	_ "errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"

	"github.com/chromedp/chromedp"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed templates/*
var templatesFS embed.FS

const (
	DefaultTimeout       = 2 * time.Minute
	DefaultMaxGoroutines = 200
	DefaultWorkers       = 10
	MaxThumbnailWidth    = 200
)

// ImageMeta holds metadata stored in DB
type ImageMeta struct {
	ID        int64
	URL       string
	Filename  string
	Thumbnail string
	Alt       string
	Title     string
	Width     int
	Height    int
	Format    string
	CrawledAt time.Time
}

// Job represents a page to crawl
type Job struct {
	URL   string
	Depth int
}

func main() {
	// CLI flags
	workerCount := flag.Int("workers", DefaultWorkers, "number of worker goroutines in pool")
	followExternal := flag.Bool("follow-external", false, "follow external links (default false)")
	enableJS := flag.Bool("enable-js", true, "enable JS rendering via chromedp for SPA pages")
	timeout := flag.Duration("timeout", DefaultTimeout, "crawling timeout, e.g. 2m")
	maxGoroutines := flag.Int("max-goroutines", DefaultMaxGoroutines, "maximum concurrent goroutines")
	imageDir := flag.String("image-dir", "images", "directory to save images and thumbnails")
	mysqlDSN := flag.String("mysql-dsn", "user:password@tcp(127.0.0.1:3306)/imagedb?parseTime=true", "MySQL DSN")
	svgRasterCmd := flag.String("svg-raster-cmd", "", "optional external command to rasterize SVGs into PNG (e.g. 'rsvg-convert -w %d -o %s %s') - provide format string with width, outpath, inputpath")
	port := flag.Int("port", 8080, "HTTP server port for search UI")
	startServer := flag.Bool("serve-only", false, "only start the web UI server (don't crawl)")
	flag.Parse()

	startURLs := flag.Args()
	if len(startURLs) == 0 && !*startServer {
		log.Fatal("provide at least one start URL as positional argument, or use -serve-only")
	}

	if *workerCount <= 0 {
		*workerCount = DefaultWorkers
	}

	if *maxGoroutines <= 0 {
		*maxGoroutines = DefaultMaxGoroutines
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	db, err := sql.Open("mysql", *mysqlDSN)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	// Ensure image directory exists
	if err := os.MkdirAll(*imageDir, 0755); err != nil {
		log.Fatalf("create image dir: %v", err)
	}

	// Start HTTP server (UI) in separate goroutine
	uiDone := make(chan struct{})
	go func() {
		if err := startHTTPServer(db, *imageDir, *port); err != nil {
			log.Printf("ui server: %v", err)
		}
		close(uiDone)
	}()

	if *startServer {
		<-uiDone
		return
	}

	// Dispatcher and worker pool
	dispatcher := NewDispatcher(*workerCount, *maxGoroutines, *followExternal, *enableJS, *imageDir, db, *svgRasterCmd)
	dispatcherCtx, dispatcherCancel := context.WithCancel(ctx)
	defer dispatcherCancel()

	go dispatcher.Run(dispatcherCtx)

	for _, u := range startURLs {
		dispatcher.Add(Job{URL: u, Depth: 0})
	}

	// wait until context expires
	<-ctx.Done()
	log.Println("main: timeout or cancelled - stopping dispatcher")
	dispatcher.Stop()
	// allow graceful shutdown of UI server for a short time
	select {
	case <-uiDone:
	case <-time.After(3 * time.Second):
	}
}

// Dispatcher orchestrates jobs and workers
type Dispatcher struct {
	workers        int
	maxGoroutines  int
	followExternal bool
	enableJS       bool
	imageDir       string
	db             *sql.DB
	svgRasterCmd   string

	jobCh   chan Job
	results chan struct{}
	quit    chan struct{}
	wg      sync.WaitGroup

	visited map[string]struct{}
	mu      sync.Mutex

	sem chan struct{} // semaphore to bound concurrent goroutines
}

func NewDispatcher(workers, maxG int, followExternal, enableJS bool, imageDir string, db *sql.DB, svgRasterCmd string) *Dispatcher {
	r := &Dispatcher{
		workers:        workers,
		maxGoroutines:  maxG,
		followExternal: followExternal,
		enableJS:       enableJS,
		imageDir:       imageDir,
		db:             db,
		svgRasterCmd:   svgRasterCmd,
		jobCh:          make(chan Job, 1000),
		results:        make(chan struct{}, 1000),
		quit:           make(chan struct{}),
		visited:        make(map[string]struct{}),
		sem:            make(chan struct{}, maxG),
	}
	return r
}

func (d *Dispatcher) Run(ctx context.Context) {
	log.Printf("dispatcher: starting with %d workers, maxGoroutines=%d\n", d.workers, d.maxGoroutines)
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.worker(ctx, i)
	}

	// Wait for cancellation
	<-ctx.Done()
	log.Println("dispatcher: context done - closing job channel")
	close(d.jobCh)
	// wait workers
	d.wg.Wait()
	log.Println("dispatcher: all workers done")
}

func (d *Dispatcher) Stop() {
	close(d.quit)
}

func (d *Dispatcher) Add(job Job) {
	// Normalize URL
	u := strings.TrimSpace(job.URL)
	if u == "" {
		return
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	job.URL = u

	d.mu.Lock()
	if _, ok := d.visited[job.URL]; ok {
		d.mu.Unlock()
		return
	}
	d.visited[job.URL] = struct{}{}
	d.mu.Unlock()

	select {
	case d.jobCh <- job:
	default:
		// job queue full; drop job (alternatively block or expand)
		log.Printf("dispatcher: job queue full, dropping %s\n", job.URL)
	}
}

func (d *Dispatcher) worker(ctx context.Context, id int) {
	defer d.wg.Done()
	log.Printf("worker %d: started\n", id)
	for job := range d.jobCh {
		select {
		case <-d.quit:
			log.Printf("worker %d: quitting\n", id)
			return
		default:
		}

		// Acquire semaphore to ensure we don't exceed global goroutine limit
		d.sem <- struct{}{}
		func(job Job) {
			defer func() { <-d.sem }()
			log.Printf("worker %d: processing %s\n", id, job.URL)
			pagesrc, baseURL, err := d.fetchPage(ctx, job.URL)
			if err != nil {
				log.Printf("worker %d: fetch %s: %v\n", id, job.URL, err)
				return
			}
			// parse page: extract links and images
			links, imgs, err := parseHTMLForLinksAndImages(bytes.NewReader(pagesrc), baseURL)
			if err != nil {
				log.Printf("worker %d: parse %s: %v\n", id, job.URL, err)
				return
			}

			// schedule links
			for _, l := range links {
				if !d.followExternal {
					if !sameSite(baseURL, l) {
						// skip externals
						continue
					}
				}
				d.Add(Job{URL: l, Depth: job.Depth + 1})
			}

			// handle images
			for _, img := range imgs {
				if err := d.processImage(ctx, img, baseURL); err != nil {
					log.Printf("worker %d: process image %s: %v\n", id, img.Src, err)
				}
			}
		}(job)
	}
	log.Printf("worker %d: stopped\n", id)
}

// fetchPage fetches page HTML. If enableJS is true it will try to render the page
// using chromedp to execute JS and return the final HTML.
func (d *Dispatcher) fetchPage(ctx context.Context, pageURL string) ([]byte, *url.URL, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, nil, err
	}

	if d.enableJS {
		// use chromedp to render page
		ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		options := chromedp.DefaultExecAllocatorOptions[:]
		allocCtx, aCancel := chromedp.NewExecAllocator(ctxt, options...)
		defer aCancel()
		cctx, cCancel := chromedp.NewContext(allocCtx)
		defer cCancel()
		var htmlContent string
		if err := chromedp.Run(cctx,
			chromedp.Navigate(pageURL),
			chromedp.Sleep(500*time.Millisecond),
			chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
		); err != nil {
			// fallback to plain HTTP fetch
			log.Printf("chromedp run failed for %s: %v - falling back to http.Get", pageURL, err)
			goto HTTPFetch
		}
		return []byte(htmlContent), u, nil
	}

HTTPFetch:
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	req.Header.Set("User-Agent", "GoImageCrawler/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return b, u, nil
}

// ImageRef represents an <img> found in a page
type ImageRef struct {
	Src   string
	Alt   string
	Title string
}

// parseHTMLForLinksAndImages parses links and image tags from HTML
func parseHTMLForLinksAndImages(r io.Reader, base *url.URL) ([]string, []ImageRef, error) {
	z := html.NewTokenizer(r)
	links := make([]string, 0)
	images := make([]ImageRef, 0)
	for {
		t := z.Next()
		switch t {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return uniqueStrings(links), images, nil
			}
			return nil, nil, z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			n := z.Token()
			if n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						if h := sanitizeURL(a.Val, base); h != "" {
							links = append(links, h)
						}
					}
				}
				continue
			}
			if n.Data == "img" {
				src := ""
				alt := ""
				title := ""
				for _, a := range n.Attr {
					switch strings.ToLower(a.Key) {
					case "src":
						src = a.Val
					case "alt":
						alt = a.Val
					case "title":
						title = a.Val
					}
				}
				if s := sanitizeURL(src, base); s != "" {
					images = append(images, ImageRef{Src: s, Alt: alt, Title: title})
				}
			}
		}
	}
}

func sanitizeURL(raw string, base *url.URL) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	// ignore data: URLs
	if strings.HasPrefix(raw, "data:") {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.Scheme == "" {
		u = base.ResolveReference(u)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	return u.String()
}

func uniqueStrings(in []string) []string {
	m := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := m[s]; !ok {
			m[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}

// sameSite reports whether u1 and u2 are from the same registrable domain (not just host).
func sameSite(base *url.URL, other string) bool {
	o, err := url.Parse(other)
	if err != nil {
		return false
	}
	bd := base.Hostname()
	oh := o.Hostname()
	bdp, _ := publicsuffix.EffectiveTLDPlusOne(bd)
	ohp, _ := publicsuffix.EffectiveTLDPlusOne(oh)
	return bdp == ohp
}

// processImage downloads image, saves file, generates thumbnail (if raster), inserts metadata to DB
func (d *Dispatcher) processImage(ctx context.Context, img ImageRef, pageBase *url.URL) error {
	u, err := url.Parse(img.Src)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u = pageBase.ResolveReference(u)
	}
	// Download image
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("User-Agent", "GoImageCrawler/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200: %d", resp.StatusCode)
	}
	// Determine filename
	fname := path.Base(u.Path)
	if fname == "" || fname == "/" || fname == "." {
		// try from URL
		fname = urlSafeFilename(u.String())
	}
	// Ensure unique filename
	unique := uniqueFilename(d.imageDir, fname)
	outfile := filepath.Join(d.imageDir, unique)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(outfile, b, 0644); err != nil {
		return err
	}
	// Try to decode image to get dimensions and type
	header := bytes.NewReader(b)
	cfg, format, err := image.DecodeConfig(header)
	thumbnailPath := ""
	width := 0
	height := 0
	if err == nil {
		width = cfg.Width
		height = cfg.Height
		if format != "svg" {
			// generate thumbnail
			thumbPath := outfile + ".thumb.png"
			if err := makeThumbnail(bytes.NewReader(b), thumbPath); err == nil {
				thumbnailPath = thumbPath
			} else {
				log.Printf("thumbnail failed: %v", err)
			}
		}
	} else {
		// Not a raster decodeable image; check if it's SVG by sniffing
		if isSVG(b) {
			format = "svg"
			// optionally rasterize using external command
			if d.svgRasterCmd != "" {
				outp := outfile + ".thumb.png"
				cmdStr := fmt.Sprintf(d.svgRasterCmd, MaxThumbnailWidth, outp, outfile)
				// Use shell to execute formatting; user must ensure command string is safe
				cmd := exec.Command("/bin/sh", "-c", cmdStr)
				if err := cmd.Run(); err == nil {
					thumbnailPath = outp
				}
			}
		} else {
			log.Printf("unknown image format for %s", outfile)
		}
	}

	// Insert into DB
	res, err := d.db.Exec(`INSERT INTO images (url, filename, thumbnail_path, alt_text, title_text, width, height, format) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, u.String(), unique, thumbnailPath, img.Alt, img.Title, width, height, format)
	if err != nil {
		return err
	}
	_ = res
	return nil
}

func urlSafeFilename(u string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9._-]`)
	return re.ReplaceAllString(u, "-")
}

func uniqueFilename(dir, fname string) string {
	base := fname
	ext := filepath.Ext(fname)
	name := strings.TrimSuffix(base, ext)
	for i := 0; ; i++ {
		candidate := fname
		if i > 0 {
			candidate = fmt.Sprintf("%s-%d%s", name, i, ext)
		}
		if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate
		}
	}
}

func makeThumbnail(r io.Reader, outpath string) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	// scale
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	if w <= MaxThumbnailWidth {
		// save as png
		f, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer f.Close()
		return png.Encode(f, img)
	}
	newW := MaxThumbnailWidth
	newH := (newW * h) / w
	newImg := image.NewRGBA(image.Rect(0, 0, newW, newH))
	// simple nearest-neighbor scaling
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := x * w / newW
			srcY := y * h / newH
			newImg.Set(x, y, img.At(srcX, srcY))
		}
	}
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, newImg)
}

func isSVG(b []byte) bool {
	s := strings.TrimSpace(string(b))
	return strings.HasPrefix(s, "<?xml") || strings.Contains(s, "<svg")
}

// startHTTPServer starts a simple web UI to search and view images
func startHTTPServer(db *sql.DB, imageDir string, port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		format := q.Get("format")
		filename := q.Get("filename")
		minw := q.Get("minw")
		minh := q.Get("minh")
		// build query
		where := []string{"1=1"}
		params := []interface{}{}
		if format != "" {
			where = append(where, "format = ?")
			params = append(params, format)
		}
		if filename != "" {
			where = append(where, "filename LIKE ?")
			params = append(params, "%"+filename+"%")
		}
		if minw != "" {
			if v, err := strconv.Atoi(minw); err == nil {
				where = append(where, "width >= ?")
				params = append(params, v)
			}
		}
		if minh != "" {
			if v, err := strconv.Atoi(minh); err == nil {
				where = append(where, "height >= ?")
				params = append(params, v)
			}
		}
		query := fmt.Sprintf("SELECT id, url, filename, thumbnail_path, alt_text, title_text, width, height, format, crawled_at FROM images WHERE %s ORDER BY crawled_at DESC LIMIT 500", strings.Join(where, " AND "))
		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, "db error", 500)
			return
		}
		defer rows.Close()
		imgs := []ImageMeta{}
		for rows.Next() {
			var im ImageMeta
			if err := rows.Scan(&im.ID, &im.URL, &im.Filename, &im.Thumbnail, &im.Alt, &im.Title, &im.Width, &im.Height, &im.Format, &im.CrawledAt); err != nil {
				log.Printf("row scan: %v", err)
				continue
			}
			imgs = append(imgs, im)
		}
		// render template
		tmplb, _ := templatesFS.ReadFile("templates/search.html")
		tmpl := string(tmplb)
		out := strings.ReplaceAll(tmpl, "{{IMAGES}}", buildImagesHTML(imgs, imageDir))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(out))
	})
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imageDir))))
	addr := fmt.Sprintf(":%d", port)
	log.Printf("http server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func buildImagesHTML(imgs []ImageMeta, dir string) string {
	var sb strings.Builder
	for _, im := range imgs {
		thumb := im.Thumbnail
		if thumb == "" {
			// use original
			thumb = filepath.Join("/images", im.Filename)
		} else {
			thumb = filepath.Join("/images", filepath.Base(thumb))
		}
		sb.WriteString("<div style='display:inline-block;margin:8px;text-align:center;width:220px'>")
		sb.WriteString(fmt.Sprintf("<a href='%s' target='_blank'><img src='%s' style='max-width:200px;display:block;margin-bottom:4px'/></a>", im.URL, thumb))
		sb.WriteString(fmt.Sprintf("<div style='font-size:12px'>%s<br/>%s %dx%d</div>", htmlEscape(im.Filename), htmlEscape(im.Format), im.Width, im.Height))
		sb.WriteString("</div>")
	}
	return sb.String()
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return replacer.Replace(s)
}
