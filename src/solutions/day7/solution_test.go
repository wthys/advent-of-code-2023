package day7

import (
	"testing"

	"github.com/wthys/advent-of-code-2023/util"
)

type (
	Result[T any, R any] struct {
		in  T
		out R
	}
)

func asHand(cards string) Hand {
	cc := Cards{}
	for _, c := range cards {
		cc = append(cc, Card(c))
	}
	return Hand{cc, 1}
}

func asCards(cards []rune) Cards {
	cc := Cards{}
	for _, c := range cards {
		cc = append(cc, Card(c))
	}
	return cc
}

func testHandTypeForPermutations(source string, expected HandType, t *testing.T) {
	util.PermutationDo(5, []rune(source), func(perm []rune) {
		hand := Hand{asCards(perm), 1}
		actual := JokerValueHandType(hand)
		if actual != expected {
			t.Fatalf("%v should be %v, got %v", hand, expected, actual)
		}
	})
}

func TestJokerValueHandType(t *testing.T) {
	testHandTypeForPermutations("ABCDE", HandType(0), t)
	testHandTypeForPermutations("ABCDD", HandType(1), t)
	testHandTypeForPermutations("AABBC", HandType(2), t)
	testHandTypeForPermutations("AAABC", HandType(3), t)
	testHandTypeForPermutations("AAABB", HandType(4), t)
	testHandTypeForPermutations("AAAAB", HandType(5), t)
	testHandTypeForPermutations("AAAAJ", HandType(6), t)
	testHandTypeForPermutations("AAABJ", HandType(5), t)
	testHandTypeForPermutations("AABCJ", HandType(3), t)
	testHandTypeForPermutations("ABCDJ", HandType(1), t)
	testHandTypeForPermutations("AAAJJ", HandType(6), t)
	testHandTypeForPermutations("AABJJ", HandType(5), t)
	testHandTypeForPermutations("ABCJJ", HandType(3), t)
	testHandTypeForPermutations("AAJJJ", HandType(6), t)
	testHandTypeForPermutations("ABJJJ", HandType(5), t)
	testHandTypeForPermutations("AJJJJ", HandType(6), t)
	testHandTypeForPermutations("JJJJJ", HandType(6), t)
}
