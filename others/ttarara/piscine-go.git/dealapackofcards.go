package piscine

import "fmt"

func DealAPackOfCards(deck []int) {
	player1 := deck[0:3]
	player2 := deck[3:6]
	player3 := deck[6:9]
	player4 := deck[9:12]

	fmt.Printf("Player 1: %d, %d, %d\n", player1[0], player1[1], player1[2])
	fmt.Printf("Player 2: %d, %d, %d\n", player2[0], player2[1], player2[2])
	fmt.Printf("Player 3: %d, %d, %d\n", player3[0], player3[1], player3[2])
	fmt.Printf("Player 4: %d, %d, %d\n", player4[0], player4[1], player4[2])
}
