package main

import (
	"fmt"
	"runtime/debug"
)

var testData = [][]int{
	{1, 2, 3, 4, 5, 5, 5},
	{-1, -2, -3, -4, -5, -5, -5},
	{}, // expected panic
	{1},
	{0},
}

// Test function which can panic
func Avg(sequence []int) int {
	sum := 0
	for _, elem := range sequence {
		sum += elem
	}
	return sum / len(sequence)
}

//CalcAvg is wrapper to run Avg safelly
func CalcAvg(sequence []int) (avg int, err error) {
	defer func() {
		panicValue := recover()
		if panicValue != nil {
			fmt.Printf("PANIC: %v\n%s", panicValue, debug.Stack())
			err = NewError(fmt.Sprintf("%v", panicValue))
		}
	}()
	avg = Avg(sequence)
	return avg, nil
}

func main() {
	for i, data := range testData {
		avg, err := CalcAvg(data)
		if err != nil {
			fmt.Printf("test %d error: %s\n", i, err.Error())
		} else {
			fmt.Printf("test %d result: avg=%d\n", i, avg)
		}
	}
	fmt.Println("all tests done")
}
