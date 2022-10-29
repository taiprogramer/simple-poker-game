package room_card

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
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

type Card struct {
	Number int
	Suit   int
}

var combinationScoreMap map[string]int = map[string]int{
	"STRAIGHT_FLUSH":  8,
	"FOUR_OF_A_KIND":  7,
	"FULL_HOUSE":      6,
	"FLUSH":           5,
	"STRAIGHT":        4,
	"THREE_OF_A_KIND": 3,
	"TWO_PAIRS":       2,
	"PAIR":            1,
	"HIGH_CARD":       0,
}

var combinationFuncMap map[string](func([]Card) bool) = map[string](func([]Card) bool){
	"STRAIGHT_FLUSH":  IsStraightFlush,
	"FOUR_OF_A_KIND":  IsFourOfAKind,
	"FULL_HOUSE":      IsFullHouse,
	"FLUSH":           IsFlush,
	"STRAIGHT":        IsStraight,
	"THREE_OF_A_KIND": IsThreeOfAKind,
	"TWO_PAIRS":       IsTwoPairs,
	"PAIR":            IsPair,
}

type Combination struct {
	Score int
}

func FindWiners(userCombinationMap map[int]Combination) map[int]Combination {
	maxScore := 0
	for _, comb := range userCombinationMap {
		if comb.Score > maxScore {
			maxScore = comb.Score
		}
	}
	winers := make(map[int]Combination)
	// find winers
	for uid, comb := range userCombinationMap {
		if comb.Score == maxScore {
			winers[uid] = comb
		}
	}
	return winers
}

func ConvertCardsStringToSortedCards(cardsStr []string) []Card {
	var cards []Card

	for _, cardStr := range cardsStr {
		number, suit := DecodeCard(cardStr)
		card := Card{
			Number: number,
			Suit:   suit,
		}
		cards = append(cards, card)
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Number < cards[j].Number
	})

	return cards
}

func IsStraight(sortedCards []Card) bool {
	length := len(sortedCards)

	for i := 0; i < length-1; i++ {
		if (sortedCards[i+1].Number - sortedCards[i].Number) != 1 {
			return false
		}
	}

	return true
}

func IsStraightFlush(sortedCards []Card) bool {
	return IsStraight(sortedCards) && IsFlush(sortedCards)
}

func IsFourOfAKind(sortedCards []Card) bool {
	var cardsMap map[int]int = make(map[int]int)

	for _, card := range sortedCards {
		n, exist := cardsMap[card.Number]
		if exist {
			n++
		} else {
			n = 1
		}
		cardsMap[card.Number] = n
	}

	for _, v := range cardsMap {
		if v == 4 {
			return true
		}
	}

	return false
}

func IsFullHouse(sortedCards []Card) bool {
	var cardsMap map[int]int = make(map[int]int)

	for _, card := range sortedCards {
		n, exist := cardsMap[card.Number]
		if exist {
			n++
		} else {
			n = 1
		}
		cardsMap[card.Number] = n
	}

	fullHouseCondition := 0
	for _, v := range cardsMap {
		if v == 3 || v == 2 {
			fullHouseCondition++
		}
	}
	return fullHouseCondition == 2
}

func IsFlush(sortedCards []Card) bool {
	length := len(sortedCards)

	for i := 0; i < length-1; i++ {
		if sortedCards[i].Suit != sortedCards[i+1].Suit {
			return false
		}
	}

	return true
}

func IsThreeOfAKind(sortedCards []Card) bool {
	var cardsMap map[int]int = make(map[int]int)

	for _, card := range sortedCards {
		n, exist := cardsMap[card.Number]
		if exist {
			n++
		} else {
			n = 1
		}
		cardsMap[card.Number] = n
	}

	for _, v := range cardsMap {
		if v == 3 {
			return true
		}
	}

	return false
}

func IsPair(sortedCards []Card) bool {
	var cardsMap map[int]int = make(map[int]int)

	for _, card := range sortedCards {
		n, exist := cardsMap[card.Number]
		if exist {
			n++
		} else {
			n = 1
		}
		cardsMap[card.Number] = n
	}

	for _, v := range cardsMap {
		if v == 2 {
			return true
		}
	}

	return false
}

func IsTwoPairs(sortedCards []Card) bool {
	var cardsMap map[int]int = make(map[int]int)

	for _, card := range sortedCards {
		n, exist := cardsMap[card.Number]
		if exist {
			n++
		} else {
			n = 1
		}
		cardsMap[card.Number] = n
	}

	numPairs := 0
	for _, v := range cardsMap {
		if v == 2 {
			numPairs++
		}
	}

	return numPairs == 2
}

func GetCombination(fiveCards []string) Combination {
	sortedCards := ConvertCardsStringToSortedCards(fiveCards[0:5])
	score := 0
	for k, f := range combinationFuncMap {
		if f(sortedCards) {
			score = combinationScoreMap[k]
		}
	}
	return Combination{
		Score: score,
	}
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func EncodeCard(number, suit int) string {
	var cards []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var suits []string = []string{"0", "1", "2", "3"}
	return cards[number-1] + suits[suit]
}

func DecodeCard(card string) (int, int) {
	var cards []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
	var suits []string = []string{"0", "1", "2", "3"}

	var number int
	var suit int
	for i, v := range cards {
		if strings.Compare(string(card[0]), v) == 0 {
			number = i + 1
			break
		}
	}
	for i, v := range suits {
		if strings.Compare(string(card[1]), v) == 0 {
			suit = i
			break
		}
	}
	return number, suit
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
