package main

import (
	"fmt"
	"os"
	"path/filepath"

	"forum/src/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s <password>", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	fmt.Println(utils.HashPassword(os.Args[1]))
}
