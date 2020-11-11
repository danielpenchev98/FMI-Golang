package main

import "fmt"

//Filter generates function which filters arguments of the function given a predicate function
//returns func(..int) []int
func Filter(p func(int) bool) func(...int) []int {
	return func(args ...int) []int {
		var result []int
		for _, item := range args {
			if p(item) {
				result = append(result, item)
			}
		}
		return result
	}
}

//Mapper generates function which transforms int into another int
//returns func(...int) []int
func Mapper(f func(int) int) func(...int) []int {
	return func(args ...int) []int {
		var result []int
		for _, item := range args {
			result = append(result, f(item))
		}
		return result
	}
}

//Reducer generates function which reduces the input to a single int
//returns func(...int) int
func Reducer(initial int, f func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, item := range args {
			initial = f(initial, item)
		}
		return initial
	}
}

//MapReducer generates function which firstly maps the arguments to new ints and reduces the modified input to single int
//returns func(...int) int
func MapReducer(initial int, mapper func(int) int, reducer func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, item := range args {
			initial = reducer(initial, mapper(item))
		}
		return initial
	}
}

func main() {
	fmt.Println("This package contains unit tests for testing the methods")
}
