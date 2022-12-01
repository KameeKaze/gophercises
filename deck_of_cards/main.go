package main

import (
	"math/rand"
	"time"
)

var (
	suits = []string{"♦", "♣", "♥", "♠"}
	ranks = []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	Deck  []Card
)

type Card struct {
	Suit string
	Rank string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	for i := range suits {
		for j := range ranks {
			Deck = append(Deck, Card{suits[i], ranks[j]})
		}
	}
}

func Shuffle(deck []Card) []Card {
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}
