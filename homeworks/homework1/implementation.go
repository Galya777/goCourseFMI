package homework1

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// subset of the GitHub user JSON fields we care about.
type User struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Followers   int    `json:"followers"`
	PublicRepos int    `json:"public_repos"`
	HTMLURL     string `json:"html_url"`
}

// subset of the GitHub repo JSON fields we care about.
type Repo struct {
	Name       string    `json:"name"`
	ForksCount int       `json:"forks_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HTMLURL    string    `json:"html_url"`
}

// Summary aggregates computed statistics for a user.
type Summary struct {
	User            User
	Repos           []Repo
	LanguageBytes   map[string]int64
	TotalForks      int
	ActivityByYear  map[int]int // year -> count (created + updated)
	TotalLanguageKB int64
}

// fetchJSON fetches the given url and decodes JSON into v.
func fetchJSON(url string, v interface{}) error {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// GitHub API requires a User-Agent header. Also Accept as json.
	req.Header.Set("User-Agent", "Go-GitHub-Client/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Basic handling for common problems.
	if resp.StatusCode == 404 {
		return errors.New("not found (404)")
	}
	if resp.StatusCode == 403 {
		// Could be rate-limited
		return fmt.Errorf("access forbidden (403). Possibly rate-limited. Response status: %s", resp.Status)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(v)
}

// readUsernames reads a file (one username per line) and returns slice of usernames.
func readUsernames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var users []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		users = append(users, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// fetchUserSummary performs all GitHub requests for a single username and computes stats.
func fetchUserSummary(username string) (*Summary, error) {
	s := &Summary{
		LanguageBytes:  make(map[string]int64),
		ActivityByYear: make(map[int]int),
	}

	// 1) fetch user
	userURL := fmt.Sprintf("https://api.github.com/users/%s", username)
	var user User
	if err := fetchJSON(userURL, &user); err != nil {
		return nil, fmt.Errorf("fetch user %s: %w", username, err)
	}
	s.User = user

	// 2) fetch repos (GitHub paginates at 30 by default; use ?per_page=100 to try to get more)
	reposURL := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", username)
	var reposRaw []struct {
		Name       string `json:"name"`
		ForksCount int    `json:"forks_count"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		HTMLURL    string `json:"html_url"`
	}
	if err := fetchJSON(reposURL, &reposRaw); err != nil {
		return nil, fmt.Errorf("fetch repos for %s: %w", username, err)
	}

	// Convert to Repo typed slice (parsing times)
	var repos []Repo
	for _, r := range reposRaw {
		created, err := time.Parse(time.RFC3339, r.CreatedAt)
		if err != nil {
			// try other format variants if any (be forgiving)
			created = time.Time{}
		}
		updated, err := time.Parse(time.RFC3339, r.UpdatedAt)
		if err != nil {
			updated = time.Time{}
		}
		repos = append(repos, Repo{
			Name:       r.Name,
			ForksCount: r.ForksCount,
			CreatedAt:  created,
			UpdatedAt:  updated,
			HTMLURL:    r.HTMLURL,
		})
		s.TotalForks += r.ForksCount
		// count activity by year using created and updated years (if valid)
		if !created.IsZero() {
			s.ActivityByYear[created.Year()]++
		}
		if !updated.IsZero() {
			s.ActivityByYear[updated.Year()]++
		}
	}
	s.Repos = repos

	// 3) for each repo fetch languages (this returns a map[string]int)
	for _, repo := range repos {
		// repo.Name is safe to use in URL
		langURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", username, repo.Name)
		var langMap map[string]int64
		if err := fetchJSON(langURL, &langMap); err != nil {
			// Non-fatal: some repos may return errors (private, removed). Skip with a warning.
			fmt.Fprintf(os.Stderr, "warning: fetch languages for %s/%s: %v\n", username, repo.Name, err)
			continue
		}
		// accumulate
		for lang, bytes := range langMap {
			s.LanguageBytes[lang] += bytes
			s.TotalLanguageKB += bytes
		}
	}

	return s, nil
}

// printSummaryTables prints the comparison table for all users and detailed tables per user.
func printSummaryTables(summaries map[string]*Summary, order []string) {
	// 1) Comparison summary table
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintln(w, "Username\tName\tFollowers\t#Repos\tTotalForks\tTotalLangKB")
	for _, uname := range order {
		s := summaries[uname]
		name := s.User.Name
		if name == "" {
			name = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\n",
			s.User.Login,
			name,
			s.User.Followers,
			len(s.Repos),
			s.TotalForks,
			s.TotalLanguageKB/1024,
		)
	}
	w.Flush()

	fmt.Println("\n=== Detailed per-user reports ===\n")

	// For each user: languages, activity by year, repos list
	for _, uname := range order {
		s := summaries[uname]
		fmt.Printf("User: %s (%s)\n", s.User.Login, nonEmpty(s.User.Name))
		fmt.Printf("Profile: %s\n", s.User.HTMLURL)
		fmt.Printf("Followers: %d | Repositories: %d | Total forks: %d\n\n", s.User.Followers, len(s.Repos), s.TotalForks)

		// Languages table
		fmt.Println("Languages (by bytes, and percentage of this user's languages):")
		if len(s.LanguageBytes) == 0 {
			fmt.Println("  No language data available.")
		} else {
			type langPair struct {
				Lang  string
				Bytes int64
			}
			var pairs []langPair
			var total int64
			for k, v := range s.LanguageBytes {
				pairs = append(pairs, langPair{Lang: k, Bytes: v})
				total += v
			}
			sort.Slice(pairs, func(i, j int) bool { return pairs[i].Bytes > pairs[j].Bytes })

			w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintln(w, "Language\tBytes\tPercent")
			for _, p := range pairs {
				percent := (float64(p.Bytes) / float64(total)) * 100
				fmt.Fprintf(w, "%s\t%d\t%0.2f%%\n", p.Lang, p.Bytes, percent)
			}
			w.Flush()
		}

		// Activity by year
		fmt.Println("\nActivity by year (counts of repo created_at + updated_at events):")
		if len(s.ActivityByYear) == 0 {
			fmt.Println("  No activity date information available.")
		} else {
			years := make([]int, 0, len(s.ActivityByYear))
			for y := range s.ActivityByYear {
				years = append(years, y)
			}
			sort.Ints(years)
			w2 := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintln(w2, "Year\tEvents")
			for _, y := range years {
				fmt.Fprintf(w2, "%d\t%d\n", y, s.ActivityByYear[y])
			}
			w2.Flush()
		}

		// Repos list (name, forks, created, updated)
		fmt.Println("\nRepositories (first 200 chars of name shown if very long):")
		if len(s.Repos) == 0 {
			fmt.Println("  No repositories found.")
		} else {
			w3 := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
			fmt.Fprintln(w3, "Name\tForks\tCreatedAt\tUpdatedAt\tURL")
			// sort repos by updated desc
			reposCopy := make([]Repo, len(s.Repos))
			copy(reposCopy, s.Repos)
			sort.Slice(reposCopy, func(i, j int) bool {
				return reposCopy[i].UpdatedAt.After(reposCopy[j].UpdatedAt)
			})
			for _, r := range reposCopy {
				created := r.CreatedAt.Format("2006-01-02")
				if r.CreatedAt.IsZero() {
					created = "-"
				}
				updated := r.UpdatedAt.Format("2006-01-02")
				if r.UpdatedAt.IsZero() {
					updated = "-"
				}
				fmt.Fprintf(w3, "%s\t%d\t%s\t%s\t%s\n", r.Name, r.ForksCount, created, updated, r.HTMLURL)
			}
			w3.Flush()
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}

func nonEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
