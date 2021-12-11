package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func findIllegal(s string) rune {
	var chunks []rune
	for _, r := range s {
		switch r {
		case '(', '[', '{', '<':
			chunks = append([]rune{r}, chunks...)
		default:
			if r == endDelimiter(chunks[0]) {
				chunks = chunks[1:]
			} else {
				fmt.Printf("Expected %s but got %s\n", string(chunks[0]), string(r))
				return r
			}
		}
	}
	return ' '
}

func endDelimiter(start rune) rune {
	switch start {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	case '<':
		return '>'
	default:
		panic(fmt.Sprintf("unexpected start delimiter %s", start))
	}
}

func score(illegal rune) int {
	switch illegal {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	case ' ':
		return 0
	default:
		panic(fmt.Sprintf("unexected end delimiter %s", illegal))
	}
}

func completionString(incomplete string) string {
	var chunks []rune
	for i, r := range incomplete {
		switch r {
		case '(', '[', '{', '<':
			chunks = append([]rune{r}, chunks...)
		default:
			if r == endDelimiter(chunks[0]) {
				chunks = chunks[1:]
			} else {
				panic(fmt.Sprintf("Corrupt line found in incompletes! %d", i))
			}
		}
	}
	remainder := ""
	for _, startDelimiter := range chunks {
		remainder += string(endDelimiter(startDelimiter))
	}
	return remainder
}

func puzzle2Score(remainder string) int {
	totalScore := 0
	for _, delim := range remainder {
		totalScore *= 5
		switch delim {
		case ')':
			totalScore += 1
		case ']':
			totalScore += 2
		case '}':
			totalScore += 3
		case '>':
			totalScore += 4
		}
		fmt.Printf(" -- %s %d\n", string(delim), totalScore)
	}
	return totalScore
}

func main() {

	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	lines := strings.Split(string(data), "\n")
	fmt.Printf("Parsing %d lines...\n", len(lines))
	cumulativeScore := 0
	var incompletes []string
	for i, line := range lines {
		illegal := findIllegal(line)
		if illegal == ' ' {
			incompletes = append(incompletes, line)
		}
		fmt.Printf("%d: %s\n", i, string(illegal))
		cumulativeScore += score(illegal)
	}
	fmt.Printf("Puzzle 1 Final score %d\n", cumulativeScore)

	var scores []int
	for i, line := range incompletes {
		c := completionString(line)
		fmt.Printf("Line %d completion %d\n", i,len(c))
		s := puzzle2Score(c)
		fmt.Printf("Completion for line %d score: %d\n%s\n%s\n\n", i, s, line, c)
		scores = append(scores, s)
	}

	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	fmt.Printf("Scores %v\nMedian %d\n", scores, scores[len(scores)/2])
}
