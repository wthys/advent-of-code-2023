package day4

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "4"
}

func (s solution) Part1(input []string) (string, error) {
	cards, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, card := range cards {
		// fmt.Printf("Card #%v => %v\n", card.id, card.Score())
		total += card.Score()
	}
	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	return solver.NotImplemented()
}

type (
	Card struct {
		id      int
		winning *set.Set[string]
		yours   *set.Set[string]
	}
)

func (card Card) String() string {
	return fmt.Sprintf("Card(id=%v, winning=%v, yours=%v)", card.id, card.winning, card.yours)
}

func (card Card) Score() int {
	matches := card.winning.Intersect(card.yours)
	switch size := matches.Len(); size {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return 1 << (size - 1)
	}
}

func parseInput(input []string) ([]Card, error) {
	cards := []Card{}
	reCard := regexp.MustCompile("^.*([0-9]+):([^\\|]*)\\|(.*)$")
	reNum := regexp.MustCompile("[0-9]+")

	for _, line := range input {
		nums := reCard.FindStringSubmatch(line)
		if nums == nil || len(nums) == 0 {
			continue
		}

		id, _ := strconv.Atoi(nums[1])
		winning := reNum.FindAllString(nums[2], -1)
		yours := reNum.FindAllString(nums[3], -1)

		card := Card{id, set.New[string](winning...), set.New[string](yours...)}
		if card.winning.Len() != len(winning) {
			return nil, fmt.Errorf("there are duplicate winning numbers")
		}
		if card.yours.Len() != len(yours) {
			return nil, fmt.Errorf("you have duplicate duplicate numbers")
		}

		cards = append(cards, card)
	}

	return cards, nil
}
