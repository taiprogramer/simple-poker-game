package room_card

import (
	"math/rand"
	"strconv"
	"time"

	cardRepo "github.com/taiprogramer/simple-poker-game/backend/repo/card"
)

// each room will have available cards
// cards encoded:
// first character is number
// second character is suit
// 13 -> A of SPADE
// 22 -> 2 of CLUB
// 30 -> 3 of DIAMOND
// 41 -> 4 of HEART
// ..
// special cases:
// T1 -> 10 of HEART
// J1 -> J of HEART
var roomCards map[int][]string = make(map[int][]string)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func ShuffleCards(roomID int) {
	var cards []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var suits []string = []string{"0", "1", "2", "3"}
	var deckOfCards []string = make([]string, 0)
	// generate new deck of cards
	for _, card := range cards {
		for _, suit := range suits {
			finalCard := card + suit
			deckOfCards = append(deckOfCards, finalCard)
		}
	}
	// shuffle deck of cards
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 52; i++ {
		index := randInt(0, 52)
		temp := deckOfCards[i]
		deckOfCards[i] = deckOfCards[index]
		deckOfCards[index] = temp
	}

	roomCards[roomID] = deckOfCards
}

func DealNextCard(roomID int) string {
	deckOfCards := roomCards[roomID]
	card := deckOfCards[len(deckOfCards)-1]
	roomCards[roomID] = deckOfCards[:len(deckOfCards)-1]
	return card
}

func GetCardID(card string) uint {
	cardNumber := map[string]int{
		"T": 10,
		"J": 11,
		"Q": 12,
		"K": 13,
	}
	number, ok := cardNumber[string(card[0])]
	if !ok {
		number, _ = strconv.Atoi(string(card[0]))
	}
	suit, _ := strconv.Atoi(string(card[1]))
	return cardRepo.FindCardIDByNumberAndSuit(number, suit)
}
