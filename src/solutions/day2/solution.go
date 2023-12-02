package day2

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

type (
	Game struct {
		id    int
		pulls []Pull
	}

	Pull struct {
		cubes map[Color]int
	}

	Color   int
	Context struct {
		lineNo int
		line   string
	}
)

const (
	Red Color = iota
	Blue
	Green
)

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "2"
}

func (context Context) String() string {
	return fmt.Sprintf("line #%v: %q", context.lineNo, context.line)
}

func parseInput(input []string) ([]Game, error) {
	games := []Game{}

	for nr, line := range input {
		context := Context{nr + 1, line}
		if len(line) == 0 {
			continue
		}
		game, err := parseGame(line, context)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func parseGame(input string, context Context) (Game, error) {
	reGame := regexp.MustCompile("^Game ([0-9]+):(.*)$")
	gameMatch := reGame.FindStringSubmatch(input)
	if len(gameMatch) == 0 {
		return Game{}, fmt.Errorf("malformed game on %v", context)
	}
	gameId, _ := strconv.Atoi(gameMatch[1])

	pulls, err := parsePulls(gameMatch[2], context)
	if err != nil {
		return Game{}, err
	}

	return Game{gameId, pulls}, nil
}

func parsePulls(input string, context Context) ([]Pull, error) {
	rePull := regexp.MustCompile("[^;]+")
	matches := rePull.FindAllString(input, -1)

	pulls := []Pull{}
	for _, pullMatch := range matches {
		pull, err := parsePull(pullMatch, context)
		if err != nil {
			return []Pull{}, err
		}
		pulls = append(pulls, pull)
	}

	return pulls, nil
}

func parsePull(input string, context Context) (Pull, error) {
	colorMap := map[string]Color{"red": Red, "blue": Blue, "green": Green}
	reCubes := regexp.MustCompile("([0-9]+) (red|green|blue)")

	cubes := map[Color]int{}
	cubeMatches := reCubes.FindAllStringSubmatch(input, -1)
	for _, cubeMatch := range cubeMatches {
		color, exist := colorMap[cubeMatch[2]]
		if !exist {
			return Pull{}, fmt.Errorf("color [%v] does not exist on %v", cubeMatch[2], context)
		}
		amount, _ := strconv.Atoi(cubeMatch[1])

		_, ok := cubes[color]
		if !ok {
			cubes[color] = 0
		}
		cubes[color] += amount
	}

	return Pull{cubes}, nil
}

func (game Game) String() string {
	return fmt.Sprintf("Game(id=%v, pulls=%v)", game.id, game.pulls)
}

func (color Color) String() string {
	switch color {
	case Red:
		return "red"
	case Blue:
		return "blue"
	case Green:
		return "green"
	default:
		return "unknown"
	}
}

func (pull Pull) String() string {
	output := "Pull("
	n := 0
	for color, amount := range pull.cubes {
		if n > 0 {
			output += ", "
		}
		output += fmt.Sprintf("%v=%v", color, amount)
		n += 1
	}
	output += ")"
	return output
}

func isPossible(pulls []Pull) bool {
	for _, pull := range pulls {
		if pull.cubes[Red] > 12 {
			return false
		}
		if pull.cubes[Green] > 13 {
			return false
		}
		if pull.cubes[Blue] > 14 {
			return false
		}
	}
	return true
}

func (s solution) Part1(input []string) (string, error) {
	games, err := parseInput(input)
	if err != nil {
		return "", err
	}

	total := 0
	for _, game := range games {
		if isPossible(game.pulls) {
			total += game.id
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	games, err := parseInput(input)
	if err != nil {
		return "", err
	}

	total := 0
	for _, game := range games {
		needed := map[Color]int{Red: 0, Blue: 0, Green: 0}
		for _, pull := range game.pulls {
			for color, amount := range needed {
				n, ok := pull.cubes[color]
				if !ok {
					n = 0
				}
				needed[color] = max(n, amount)
			}
		}

		power := 1
		for _, amount := range needed {
			power *= amount
		}
		total += power
	}

	return solver.Solved(total)
}
