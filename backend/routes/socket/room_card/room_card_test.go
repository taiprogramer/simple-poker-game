package room_card

import "testing"

const expectedButActual = "Expected: %v but Actual %v"

func CheckExpect(expect, actual bool, t *testing.T) {
	if expect != actual {
		t.Fatalf(expectedButActual, expect, actual)
	}
}

func TestIsStraight(t *testing.T) {
	notAStraightCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "J1", "J2", "31"})
	straightCards := ConvertCardsStringToSortedCards([]string{
		"21", "32", "41", "53", "63"})

	t.Run("must return false if combination is not straight",
		func(t *testing.T) {
			CheckExpect(false, IsStraight(notAStraightCards), t)
		})

	t.Run("must return true if combination is straight",
		func(t *testing.T) {
			CheckExpect(true, IsStraight(straightCards), t)
		})
}

func TestIsPair(t *testing.T) {
	notContainsPairCards := ConvertCardsStringToSortedCards([]string{
		"21", "32", "71", "91", "81"})
	pairCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "41", "51", "91"})

	t.Run("must return true if combination is pair",
		func(t *testing.T) {
			CheckExpect(true, IsPair(pairCards), t)
		})

	t.Run("must return false if combination is not pair",
		func(t *testing.T) {
			CheckExpect(false, IsPair(notContainsPairCards), t)
		})
}

func TestIsStraightFlush(t *testing.T) {
	notStraightFlushCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "71", "82", "91"})
	straightFlushCards := ConvertCardsStringToSortedCards([]string{
		"21", "31", "41", "51", "61"})
	straightButNotFlushCards := ConvertCardsStringToSortedCards([]string{
		"21", "32", "41", "51", "61"})

	t.Run("must return true if combination is straight flush",
		func(t *testing.T) {
			CheckExpect(true, IsStraightFlush(straightFlushCards), t)
		})
	t.Run("must return false if combination is not straight flush",
		func(t *testing.T) {
			CheckExpect(false, IsStraightFlush(notStraightFlushCards), t)
		})
	t.Run("must return false if combination is straight but not flush",
		func(t *testing.T) {
			CheckExpect(false,
				IsStraightFlush(straightButNotFlushCards), t)
		})
}

func TestIsFourOfAKind(t *testing.T) {
	fourKindCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "23", "20", "31"})
	notFourKindCards := ConvertCardsStringToSortedCards([]string{
		"21", "32", "23", "20", "31"})

	t.Run("must return true if combination is four of a kind",
		func(t *testing.T) {
			CheckExpect(true, IsFourOfAKind(fourKindCards), t)
		})
	t.Run("must return false if combination is not four of a kind",
		func(t *testing.T) {
			CheckExpect(false, IsFourOfAKind(notFourKindCards), t)
		})
}

func TestIsFullHouse(t *testing.T) {
	fullHouseCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "23", "71", "72"})
	notFullHouseCards := ConvertCardsStringToSortedCards([]string{
		"21", "22", "23", "71", "81"})
	t.Run("must return true if combination is full house",
		func(t *testing.T) {
			CheckExpect(true, IsFullHouse(fullHouseCards), t)
		})
	t.Run("must return false if combination is not full house",
		func(t *testing.T) {
			CheckExpect(false, IsFullHouse(notFullHouseCards), t)
		})
}

func TestIsFlush(t *testing.T) {
	flushCards := ConvertCardsStringToSortedCards([]string{
		"71", "11", "T1", "91", "31"})
	notFlushCards := ConvertCardsStringToSortedCards([]string{
		"71", "12", "T1", "91", "31"})
	t.Run("must return true if combination is flush",
		func(t *testing.T) {
			CheckExpect(true, IsFlush(flushCards), t)
		})
	t.Run("must return false if combination is not flush",
		func(t *testing.T) {
			CheckExpect(false, IsFlush(notFlushCards), t)
		})
}

func TestIsThreeOfAKind(t *testing.T) {
	threeOfKindCards := ConvertCardsStringToSortedCards([]string{
		"71", "71", "73", "12", "11"})
	notThreeOfKindCards := ConvertCardsStringToSortedCards([]string{
		"71", "71", "73", "12", "11"})
	t.Run("must return true if combination is three of a kind",
		func(t *testing.T) {
			CheckExpect(true, IsThreeOfAKind(threeOfKindCards), t)
		})
	t.Run("must return false if combination is not three of a kind",
		func(t *testing.T) {
			CheckExpect(true, IsThreeOfAKind(notThreeOfKindCards), t)
		})
}

func TestIsTwoPairs(t *testing.T) {
	twoPairsCards := ConvertCardsStringToSortedCards([]string{
		"71", "72", "91", "91", "83"})
	onePairsCards := ConvertCardsStringToSortedCards([]string{
		"71", "72", "91", "61", "83"})
	t.Run("must return true if combination is two pairs",
		func(t *testing.T) {
			CheckExpect(true, IsTwoPairs(twoPairsCards), t)
		})
	t.Run("must return false if combination is not two pairs",
		func(t *testing.T) {
			CheckExpect(false, IsTwoPairs(onePairsCards), t)
		})
}
