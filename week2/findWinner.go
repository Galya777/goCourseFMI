package week2

import "fmt"

func findWinner(n, m int) int {
	people := make([]int, n)
	for i := 0; i < n; i++ {
		people[i] = i + 1
	}
	index := 0

	for len(people) > 1 {
		index = (index + m - 1) % len(people)
		people = append(people, people[index+1:]...)
	}
	return people[0]
}
func main() {
	var n, m int
	_, _ = fmt.Scan(&n, &m)
	fmt.Println(findWinner(n, m))
}
