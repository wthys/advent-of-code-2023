package day15

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
	return "15"
}

func (s solution) Part1(input []string) (string, error) {
	instructions, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, instr := range instructions {
		total += instr.Hash()
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	instructions, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	storage := LensStorage{}
	for _, instr := range instructions {
		lens := instr.ToLens()
		location := Hash(lens.label)
		lenses, ok := storage[location]
		if !ok {
			lenses = Lenses{}
		}
		if lens.focal > 0 {
			lenses = lenses.Store(lens)
		} else {
			lenses = lenses.Remove(lens)
		}
		storage[location] = lenses
	}

	total := 0
	for box, lenses := range storage {
		for n, lens := range lenses {
			total += (box + 1) * (n + 1) * lens.focal
		}
	}

	return solver.Solved(total)
}

type (
	Instruction  string
	Instructions []Instruction
	Lens         struct {
		label string
		focal int
	}
	Lenses      []Lens
	LensStorage map[int]Lenses
)

var (
	reInstr = regexp.MustCompile("^([a-z]+)(-|=[1-9])$")
)

func (instr Instruction) Hash() int {
	return Hash(string(instr))
}

func Hash(instr string) int {
	hsh := 0
	for _, char := range instr {
		hsh = (17 * (hsh + int(char))) % 256
	}
	return hsh
}

func (lenses Lenses) Remove(lens Lens) Lenses {
	newLenses := Lenses{}

	for _, current := range lenses {
		if current.label != lens.label {
			newLenses = append(newLenses, current)
		}
	}

	return newLenses
}

func (lenses Lenses) Store(lens Lens) Lenses {
	newLenses := Lenses{}

	replaced := false
	for _, current := range lenses {
		if current.label == lens.label {
			newLenses = append(newLenses, lens)
			replaced = true
		} else {
			newLenses = append(newLenses, current)
		}
	}

	if !replaced {
		newLenses = append(newLenses, lens)
	}

	return newLenses
}

func (instr Instruction) ToLens() Lens {
	matches := reInstr.FindStringSubmatch(string(instr))
	if len(matches) == 0 {
		panic(fmt.Sprintf("could not find label in %q", string(instr)))
	}

	label := matches[1]

	focal := 0
	switch matches[2] {
	case "-":
		focal = 0
	default:
		focal = int(matches[2][1] - '0')
	}

	return Lens{label, focal}
}

func ParseInput(input []string) (Instructions, error) {
	instructions := Instructions{}
	reSep := regexp.MustCompile("[^,]+")

	for _, line := range input {
		matches := reSep.FindAllString(line, -1)

		for _, instr := range matches {
			instructions = append(instructions, Instruction(instr))
		}
	}

	if len(instructions) == 0 {
		return Instructions{}, fmt.Errorf("no instructions found")
	}
	return instructions, nil
}
