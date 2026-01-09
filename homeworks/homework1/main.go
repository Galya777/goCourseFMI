package homework1

import (
	"fmt"
	"os"
)

func Main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <users-file>\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	usernames, err := readUsernames(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading usernames file: %v\n", err)
		os.Exit(2)
	}
	if len(usernames) == 0 {
		fmt.Fprintf(os.Stderr, "no usernames found in file\n")
		os.Exit(3)
	}

	// For each username, fetch and compute summary.
	summaries := make(map[string]*Summary)
	var order []string
	for _, u := range usernames {
		fmt.Fprintf(os.Stderr, "processing %s...\n", u)
		s, err := fetchUserSummary(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error fetching data for %s: %v\n", u, err)
			// continue to next user (so we can compare available ones)
			continue
		}
		summaries[u] = s
		order = append(order, u)
		// small delay could be added to be polite (not added by default)
	}

	if len(summaries) == 0 {
		fmt.Fprintf(os.Stderr, "no user data fetched successfully\n")
		os.Exit(4)
	}

	printSummaryTables(summaries, order)
}
