package helpers

// function that implements Levenshtein distance (edit distance) that measures the difference between to strings
func MatchToSuggest(a string, b string) int {
	la := len(a)
	lb := len(b)

	dp := make([][]int, la+1)
	for i := 0; i <= la; i++ {
		dp[i] = make([]int, lb+1)
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = Min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)

		}

	}
	return dp[la][lb]
}

func Min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}

	if b < c {
		return b
	}
	return c
}
