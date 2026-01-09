package week4

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Edge struct {
	To     int
	Weight float64
}

type Graph struct {
	Cities map[int]string
	Adj    map[int][]Edge
}

func Dijkstra(g *Graph, start, end int) ([]int, float64) {
	dist := map[int]float64{}
	prev := map[int]int{}
	visited := map[int]bool{}

	for id := range g.Cities {
		dist[id] = math.Inf(1)
	}
	dist[start] = 0

	for {
		min := math.Inf(1)
		u := -1
		for id := range dist {
			if !visited[id] && dist[id] < min {
				min = dist[id]
				u = id
			}
		}
		if u == -1 || u == end {
			break
		}
		visited[u] = true

		for _, e := range g.Adj[u] {
			if dist[u]+e.Weight < dist[e.To] {
				dist[e.To] = dist[u] + e.Weight
				prev[e.To] = u
			}
		}
	}

	path := []int{}
	for at := end; at != 0; at = prev[at] {
		path = append([]int{at}, path...)
		if at == start {
			break
		}
	}
	return path, dist[end]
}

func LoadGraph(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	g := &Graph{
		Cities: map[int]string{},
		Adj:    map[int][]Edge{},
	}

	scanner := bufio.NewScanner(file)
	readCities := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			readCities = false
			continue
		}

		p := strings.Split(line, ",")
		if readCities {
			id, _ := strconv.Atoi(p[0])
			g.Cities[id] = p[1]
		} else {
			a, _ := strconv.Atoi(p[0])
			b, _ := strconv.Atoi(p[1])
			d, _ := strconv.ParseFloat(strings.ReplaceAll(p[2], ",", "."), 64)

			g.Adj[a] = append(g.Adj[a], Edge{b, d})
			g.Adj[b] = append(g.Adj[b], Edge{a, d})
		}
	}
	return g, nil
}

func handler(g *Graph) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var result string
		var errMsg string

		if r.Method == "POST" {
			start := r.FormValue("start")
			end := r.FormValue("end")

			if start == "" || end == "" {
				errMsg = "Моля изберете начален и краен пункт!"
			} else {
				s, _ := strconv.Atoi(start)
				e, _ := strconv.Atoi(end)
				path, dist := Dijkstra(g, s, e)

				result = "<h3>Най-кратък маршрут:</h3><ul>"
				for i := 0; i < len(path); i++ {
					result += "<li>" + g.Cities[path[i]] + "</li>"
				}
				result += "</ul>"
				result += fmt.Sprintf("<b>Общо разстояние:</b> %.2f km", dist)
			}
		}

		fmt.Fprint(w, `
		<html><body>
		<h2>Оптимизация на маршрут</h2>
		<form method="POST">
			Начало:
			<select name="start">
				<option value="">--избор--</option>`)

		for id, name := range g.Cities {
			fmt.Fprintf(w, `<option value="%d">%s</option>`, id, name)
		}

		fmt.Fprint(w, `
			</select>
			Край:
			<select name="end">
				<option value="">--избор--</option>`)

		for id, name := range g.Cities {
			fmt.Fprintf(w, `<option value="%d">%s</option>`, id, name)
		}

		fmt.Fprint(w, `
			</select>
			<br><br>
			<input type="submit" value="Намери маршрут">
		</form>
		<p style="color:red;">`+errMsg+`</p>
		`+result+`
		</body></html>`)
	}
}

func main() {
	file := flag.String("file", "map.txt", "input file")
	flag.Parse()

	graph, err := LoadGraph(*file)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler(graph))
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
