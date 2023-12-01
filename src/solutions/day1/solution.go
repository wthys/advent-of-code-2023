package day1

import (
    "fmt"
    "regexp"
    "strconv"
    "github.com/wthys/advent-of-code-2023/solver"
)


type solution struct {}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "1"
}

func (s solution) Part1(input []string) (string, error) {
    re := regexp.MustCompile("[0-9]")

    total := 0

    for _, line := range input {
        matches := re.FindAllString(line, -1)
        
        if (len(matches) > 0) {
            cand, _ := strconv.Atoi(matches[0] + matches[len(matches)-1])
            total = total + cand
        }
    }


    return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
    return solver.NotImplemented()
}
