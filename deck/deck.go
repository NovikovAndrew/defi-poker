package deck

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades:
		return "Spades"
	case Harts:
		return "Harts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	default:
		return ""
	}
}

func (s Suit) suitToUnicode() string {
	switch s {
	case Spades:
		return "♠"
	case Harts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		return ""
	}
}

const (
	Spades Suit = iota
	Harts
	Diamonds
	Clubs
)

type Card struct {
	suit  Suit
	value int
}

func NewCard(s Suit, v int) (Card, error) {
	if v > 13 {
		return Card{}, errors.New("the value can not be great 13")
	}

	return Card{
		suit:  s,
		value: v,
	}, nil
}

func (c Card) String() string {
	value := strconv.Itoa(c.value)
	if c.value == 1 {
		value = "Ace"
	}

	return fmt.Sprintf("%s of %s %s", value, c.suit.String(), c.suit.suitToUnicode())
}

type Deck [52]Card

func NewDeck() (Deck, error) {
	var d = [52]Card{}
	nSuits := 4
	nCards := 13
	x := 0

	for i := 0; i < nSuits; i++ {
		for j := 0; j < nCards; j++ {
			c, err := NewCard(Suit(i), j+1)
			if err != nil {
				return Deck{}, err
			}

			d[x] = c
			x++
		}
	}

	return shuffle(d), nil
}

func shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)
		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}

	return d
}
