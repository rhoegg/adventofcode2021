package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Fish struct {
	age, count int
}

func CountFish(fish []Fish) int {
	var count int
	for _, f := range fish {
		count += f.count
	}
	return count
}

func ageOneDay(fish []Fish)  []Fish {
	var output []Fish
	newborns := Fish{age: 8, count: 0}
	for _, f := range fish {
		if f.age == 0 {
			newborns.count += f.count
			f.age = 7
		}
		f.age -= 1
		output = append(output, f)
	}
	return append(output, newborns)
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	inputFish := strings.Split(string(data), ",")
	fmt.Printf("fish: %d\n", len(inputFish))

	var fish []Fish
	for _, f := range inputFish {
		i, _ := strconv.Atoi(f)
		fish = append(fish, Fish{age: i, count: 1})
	}

	for i := 1; i <= 256; i++ {
		fish = ageOneDay(fish)
		fmt.Printf("After Day %d: %d\n", i, CountFish(fish))
	}
	fmt.Printf("%v\n", fish)
}
