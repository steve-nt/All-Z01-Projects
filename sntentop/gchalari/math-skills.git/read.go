package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func readNumbers() []int {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a file name")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var numbers []int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		num, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}

		numbers = append(numbers, num)
	}

	return numbers
}
