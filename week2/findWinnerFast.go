package week2

func findWinnerFast(n int, m int) int {
	winner := 0

	for i := 2; i < n; i++ {
		winner = (winner + m) % i
	}

	return winner + 1
}
