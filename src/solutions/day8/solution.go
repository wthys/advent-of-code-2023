package day8

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util"
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
	nodes := desertMap.Nodes()
	slices.Sort(nodes)
	current := nodes[0]
	target := nodes[len(nodes)-1]
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
	desertMap, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	reStart := regexp.MustCompile(`A$`)
	reEnd := regexp.MustCompile(`Z$`)

	endStates := []int{}
	diffs := []int{}
	current := Nodes{}
	desertMap.Nodes().ForEach(func(node Node) {
		if reStart.MatchString(string(node)) {
			current = append(current, node)
			endStates = append(endStates, 0)
			diffs = append(diffs, 0)
		}
	})

	step := 0
	limit := 50000
	stopped := true
	for step < limit {
		next := Nodes{}
		ended := true
		for idx, node := range current {
			nn, ok := desertMap.Next(node, step)
			if !ok {
				return solver.Error(fmt.Errorf("could not find next node for %v", node))
			}
			next = append(next, nn)
			if reEnd.MatchString(string(nn)) {
				diff := step - endStates[idx]
				diffs[idx] = diff
				limit = max(limit, 5*diff)
				endStates[idx] = step
			} else {
				ended = false
			}
		}
		current = next
		step += 1
		if ended {
			stopped = false
			break
		}
	}

	if stopped {
		switch len(diffs) {
		case 1:
			step = diffs[0]
		case 2:
			step = util.LCM(diffs[0], diffs[1])
		default:
			step = util.LCM(diffs[0], diffs[1], diffs[2:]...)
		}
	}

	return solver.Solved(step)
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

func (m Map) Nodes() Nodes {
	nodes := Nodes{}
	for key, _ := range m.nodes {
		nodes = append(nodes, key)
	}
	return nodes
}

func (nodes Nodes) ForEach(forEach func(Node)) {
	for _, node := range nodes {
		forEach(node)
	}
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
	reNodeMap := regexp.MustCompile(`^\s*([0-9A-Z]+)\s*=\s*[(]\s*([0-9A-Z]+)\s*,\s*([0-9A-Z]+)\s*[)]\s*$`)

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
