package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

const (
	windowSize = 20
	outlierCut = 1000
)

func median(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]int{}, values...)
	sort.Ints(sorted)
	n := len(sorted)
	if n%2 == 1 {
		return float64(sorted[n/2])
	}
	return float64(sorted[n/2-1]+sorted[n/2]) / 2
}

func main() {
	idealRanges := []int{3, 13, 28, 34, 41, 59, 84, 94, 123, 145, 177, 228, 320, 533, 1600}
	min := math.MaxInt
	max := math.MinInt

	scanner := bufio.NewScanner(os.Stdin)
	history := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("103 196")
			continue
		}

		med := median(history)

		if len(history) == 0 || math.Abs(float64(num)-med) < outlierCut {
			history = append(history, num)
			if min > num {
				min = num
			}
			if max < num {
				max = num
			}
		}

		if len(history) > windowSize {
			history = history[1:]
		}

		diff := max - min
		var rangeWidth int
		for _, val := range idealRanges {
			if diff > val {
				rangeWidth = val
			} else {
				break
			}
		}

		var lower, upper int
		if rangeWidth == 0 {
			fmt.Printf("%d %d\n", min, max)
			continue
		} else if rangeWidth%2 == 1 {
			lower = (min+max)/2 - rangeWidth/2
			upper = (min+max)/2 + rangeWidth/2
		} else {
			lower = (min+max)/2 - rangeWidth/2
			upper = (min+max)/2 + rangeWidth/2 - 1
		}
		fmt.Printf("%d %d\n", lower, upper)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
