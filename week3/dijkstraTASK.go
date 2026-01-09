package week3

import (
	"bufio"
	"fmt"
	"math"
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

func NewGraph() *Graph {
	return &Graph{
		Cities: make(map[int]string),
		Adj:    make(map[int][]Edge),
	}
}
func LoadGraph(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := NewGraph()
	scanner := bufio.NewScanner(file)
	readingCities := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			readingCities = false
			continue
		}

		parts := strings.Split(line, ",")

		if readingCities {
			id, _ := strconv.Atoi(parts[0])
			graph.Cities[id] = parts[1]
		} else {
			from, _ := strconv.Atoi(parts[0])
			to, _ := strconv.Atoi(parts[1])
			dist, _ := strconv.ParseFloat(strings.ReplaceAll(parts[2], ",", "."), 64)

			graph.Adj[from] = append(graph.Adj[from], Edge{to, dist})
			graph.Adj[to] = append(graph.Adj[to], Edge{from, dist}) // двупосочно
		}
	}

	return graph, nil
}

func Dijkstra(g *Graph, start, end int) ([]int, float64) {
	dist := make(map[int]float64)
	prev := make(map[int]int)
	visited := make(map[int]bool)

	for city := range g.Cities {
		dist[city] = math.Inf(1)
	}
	dist[start] = 0

	for {
		minDist := math.Inf(1)
		current := -1

		for city := range g.Cities {
			if !visited[city] && dist[city] < minDist {
				minDist = dist[city]
				current = city
			}
		}

		if current == -1 || current == end {
			break
		}

		visited[current] = true

		for _, edge := range g.Adj[current] {
			newDist := dist[current] + edge.Weight
			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = current
			}
		}
	}

	// възстановяване на пътя
	path := []int{}
	for at := end; at != 0; at = prev[at] {
		path = append([]int{at}, path...)
		if at == start {
			break
		}
	}

	return path, dist[end]
}

func main() {
	graph, err := LoadGraph("map.txt")
	if err != nil {
		fmt.Println("Грешка:", err)
		return
	}

	var start, end int
	fmt.Print("Начален град (ID): ")
	fmt.Scan(&start)
	fmt.Print("Краен град (ID): ")
	fmt.Scan(&end)

	path, distance := Dijkstra(graph, start, end)

	fmt.Println("\nНай-кратък маршрут:")
	for i := 0; i < len(path); i++ {
		fmt.Print(graph.Cities[path[i]])
		if i < len(path)-1 {
			fmt.Print(" -> ")
		}
	}

	fmt.Printf("\nОбщо разстояние: %.2f km\n", distance)
}
