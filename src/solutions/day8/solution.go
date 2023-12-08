package day8

import (
	"fmt"
	"regexp"

	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "8"
}

func (s solution) Part1(input []string) (string, error) {
	desertMap, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	step := 0
	current := Node("AAA")
	target := Node("ZZZ")
	for current != target {
		next, ok := desertMap.Next(current, step)
		if !ok {
			return solver.Error(fmt.Errorf("could not find next node for %v", current))
		}
		current = next
		step += 1
	}

	return solver.Solved(step)
}

func (s solution) Part2(input []string) (string, error) {
	return solver.NotImplemented()
}

type (
	Map struct {
		instructions string
		nodes        NodeMap
	}

	NodeMap map[Node]Nodes
	Node    string
	Nodes   []Node
)

func (m Map) Next(node Node, step int) (Node, bool) {
	paths, ok := m.nodes[node]
	if !ok {
		return Node(""), false
	}

	instr := m.instructions[step%len(m.instructions)]

	next := Node("")
	if instr == 'L' {
		next = paths[0]
	} else {
		next = paths[1]
	}
	// fmt.Printf("%v: %v -%v-> %v\n", step+1, node, string(instr), next)
	return next, true
}

func ParseInput(input []string) (Map, error) {
	if len(input) < 3 {
		return Map{}, fmt.Errorf("not enough input")
	}
	instructions := input[0]
	if len(instructions) == 0 {
		return Map{}, fmt.Errorf("no instructions found")
	}
	reInstr := regexp.MustCompile(`^[LR]+$`)
	if !reInstr.MatchString(instructions) {
		return Map{}, fmt.Errorf("invalid instructions found")
	}

	nodes := NodeMap{}
	reNodeMap := regexp.MustCompile(`^\s*([A-Z]+)\s*=\s*[(]\s*([A-Z]+)\s*,\s*([A-Z]+)\s*[)]\s*$`)

	for lineNr, line := range input[2:] {
		matches := reNodeMap.FindStringSubmatch(line)

		if len(matches) == 0 {
			continue
		}

		if len(matches) < 4 {
			return Map{}, fmt.Errorf("could not parse line #%v : %q", lineNr+3, line)
		}

		nodes[Node(matches[1])] = Nodes{Node(matches[2]), Node(matches[3])}
	}

	return Map{instructions, nodes}, nil
}
