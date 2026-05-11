package main

import (
	"fmt"
	"go-reloaded/core" // Adjust this import path based on your project structure
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Display the menu for text selection
	fmt.Println("Choose a text scenario to process:")
	fmt.Println("1) If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?")
	fmt.Println("2) I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure")
	fmt.Println("3) Don not be sad ,because sad backwards is das . And das not good")
	fmt.Println("4) harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '")

	// Read the user's choice
	var choice int
	fmt.Print("Enter your choice (1-4): ")
	_, err := fmt.Scan(&choice)
	if err != nil || choice < 1 || choice > 4 {
		log.Fatal("Invalid choice. Please run the program again.")
	}

	// Determine which scenario to use based on user input
	var sampleText string
	switch choice {
	case 1:
		sampleText = "If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?"
	case 2:
		sampleText = "I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure"
	case 3:
		sampleText = "Don not be sad ,because sad backwards is das . And das not good"
	case 4:
		sampleText = "harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '"
	}

	// Write the chosen scenario to sample.txt
	err = ioutil.WriteFile("sample.txt", []byte(sampleText), 0644)
	if err != nil {
		log.Fatal("Error writing to sample.txt:", err)
	}

	fmt.Println("Successfully wrote to sample.txt")

	// Check if the correct number of arguments is provided
	if len(os.Args) != 3 {
		log.Fatal("Invalid argument count. Expect 2 arguments corresponding to input/output file names.")
	}

	// Call the core.Run() function to process the sample.txt
	core.Run()
}
