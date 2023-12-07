package day7

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "7"
}

func (s solution) Part1(input []string) (string, error) {
	hands, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	types := map[HandType]Hands{}

	for _, hand := range hands {
		handType := FaceValueHandType(hand)
		more, ok := types[handType]
		if !ok {
			more = Hands{}
		}
		more = append(more, hand)
		types[handType] = more
	}

	const ranking = "23456789TJQKA"

	for _, hands := range types {
		slices.SortFunc(hands, handComparator(ranking))
	}

	rank := 1
	total := 0
	ttt := []HandType{0, 1, 2, 3, 4, 5, 6}
	for _, t := range ttt {
		for _, hand := range types[t] {
			total += rank * hand.bid
			rank += 1
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	hands, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	types := map[HandType]Hands{}

	for _, hand := range hands {
		handType := JokerValueHandType(hand)
		more, ok := types[handType]
		if !ok {
			more = Hands{}
		}
		more = append(more, hand)
		types[handType] = more
	}

	const ranking = "J23456789TQKA"

	for _, hands := range types {
		slices.SortFunc(hands, handComparator(ranking))
	}

	rank := 1
	total := 0
	ttt := []HandType{0, 1, 2, 3, 4, 5, 6}
	for _, t := range ttt {
		for _, hand := range types[t] {
			total += rank * hand.bid
			rank += 1
		}
	}

	return solver.Solved(total)
}

type (
	Card  rune
	Cards []Card
	Hand  struct {
		cards Cards
		bid   int
	}
	Hands     []Hand
	CardCount map[Card]int

	HandType int
)

func (cards Cards) Count() CardCount {
	counts := CardCount{}
	for _, card := range cards {
		counts[card] += 1
	}
	return counts
}

func (cc CardCount) GetN(n int) Cards {
	cards := Cards{}
	for card, count := range cc {
		if count == n {
			cards = append(cards, card)
		}
	}
	return cards
}

func (cc CardCount) Occurences(card Card) int {
	amount, ok := cc[card]
	if !ok {
		return 0
	}
	return amount
}

func handComparator(ranking string) func(a, b Hand) int {
	return func(a, b Hand) int {
		for idx := 0; idx < len(a.cards); idx++ {
			cmp := compareCards(ranking, a.cards[idx], b.cards[idx])
			if cmp != 0 {
				return cmp
			}
		}
		return 0
	}
}
func compareCards(ranking string, l, r Card) int {
	if l == r {
		return 0
	}
	lr := strings.IndexRune(ranking, rune(l))
	rr := strings.IndexRune(ranking, rune(r))
	return util.Sign(lr - rr)
}

func FaceValueHandType(hand Hand) HandType {
	count := hand.cards.Count()
	if len(count.GetN(5)) > 0 {
		return 6
	}

	if len(count.GetN(4)) > 0 {
		return 5
	}

	if len(count.GetN(3)) > 0 {
		if len(count.GetN(2)) > 0 {
			return 4
		}
		return 3
	}

	return HandType(len(count.GetN(2)))
}

func JokerValueHandType(hand Hand) HandType {
	count := hand.cards.Count()
	jokers := count.Occurences('J')
	count['J'] = 0
	if len(count.GetN(5)) > 0 {
		return 6
	}

	if faces := count.GetN(4); len(faces) > 0 {
		if jokers > 0 {
			return 6
		}
		return 5
	}

	if faces := count.GetN(3); len(faces) > 0 {
		if jokers > 0 {
			return HandType(4 + jokers)
		} else {
			if len(count.GetN(2)) > 0 {
				return 4
			}
			return 3
		}
	}

	faces := count.GetN(2)
	switch len(faces) {
	case 0:
		switch jokers {
		case 0, 1:
			return HandType(jokers)
		case 2:
			return 3
		default:
			return HandType(min(6, jokers+2))
		}
	case 1:
		return HandType(min(2*jokers+1, 6))
	case 2:
		if jokers > 0 {
			return 4
		}
		return 2
	default:
		panic(fmt.Sprintf("cannot have more than 2 pairs, got %v", len(faces)))
		return 0
	}
}

var TYPES = []string{"high", "1pair", "2pair", "3kind", "fullh", "4kind", "5kind"}

func (ht HandType) String() string {
	return TYPES[ht]
}

func (this Hand) String() string {
	return fmt.Sprintf("[%v %v]", this.cards, this.bid)
}

func (this Cards) String() string {
	out := ""
	for _, card := range this {
		out += string(card)
	}
	return out
}

func parseInput(input []string) (Hands, error) {
	hands := Hands{}
	reHand := regexp.MustCompile("^([2-9TJQKA]{5}) ([0-9]+)$")

	for _, line := range input {
		matches := reHand.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		cs := []Card{}
		for _, card := range matches[1] {
			cs = append(cs, Card(card))
		}
		cards := Cards(cs)

		bid, _ := strconv.Atoi(matches[2])
		hands = append(hands, Hand{cards, bid})
	}

	return hands, nil
}
