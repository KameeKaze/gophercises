package main

import "fmt"

var (
	suits = []string{"♦", "♣", "♥", "♠"}
	ranks = []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	deck  []Card
)

type Card struct {
	Suit string
	Rank string
}

func main() {
	for i := range suits {
		for j := range ranks {
			deck = append(deck, Card{suits[i], ranks[j]})
			fmt.Println(suits[i], ranks[j])
		}

	}

}
