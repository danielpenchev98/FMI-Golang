package main

import (
	"testing"
)

func TestFilter(t *testing.T) {
	table := []struct {
		input  []int
		output []int
		descr  string
	}{
		{input: []int{1, 2, 3, 4, 5}, output: []int{1, 3, 5}, descr: "Filter some elements"},
		{input: []int{1, 3, 5}, output: []int{1, 3, 5}, descr: "No elements filtered"},
		{input: []int{}, output: []int{}, descr: "Empty array"},
	}

	odds := Filter(func(x int) bool { return x%2 == 1 })

	for _, testCase := range table {
		actualOutput := odds(testCase.input...)
		expectedOutput := testCase.output

		if len(actualOutput) != len(expectedOutput) {
			t.Errorf("Test Case [%s].The length actual output differs from the length of the expected one", testCase.descr)
		}

		for idx := range actualOutput {
			if actualOutput[idx] != expectedOutput[idx] {
				t.Errorf("Test Case [%s]. Actual output [%v], but wanted [%v]", testCase.descr, actualOutput, expectedOutput)
			}
		}
	}
}

func TestMapper(t *testing.T) {
	table := []struct {
		input  []int
		output []int
		descr  string
	}{
		{input: []int{1, 2, 3}, output: []int{2, 4, 6}, descr: "Apply function on all elements"},
		{input: []int{}, output: []int{}, descr: "Empty array"},
	}

	double := Mapper(func(a int) int { return 2 * a })

	for _, testCase := range table {
		actualOutput := double(testCase.input...)
		expectedOutput := testCase.output

		if len(actualOutput) != len(expectedOutput) {
			t.Errorf("Test Case [%s].The length actual output differs from the length of the expected one", testCase.descr)
		}

		for idx := range actualOutput {
			if actualOutput[idx] != expectedOutput[idx] {
				t.Errorf("Test Case [%s]. Actual output [%v], but wanted [%v]", testCase.descr, actualOutput, expectedOutput)
			}
		}
	}
}

func TestReducer(t *testing.T) {
	sum := Reducer(0, func(a, b int) int { return a + b })
	testTable := []struct {
		input  []int
		output int
	}{
		{input: []int{1, 2, 3}, output: 6},
		{input: []int{5}, output: 11},
		{input: []int{100, 101, 102}, output: 314},
	}

	for i := 0; i < len(testTable); i++ {
		testCase := testTable[i]
		output := sum(testCase.input...)
		if testCase.output != output {
			t.Errorf("Call %d, expected output %v, actual output %v", i+1, testCase.output, output)
		}
	}
}

func TestMapReducer(t *testing.T) {
	powerSum := MapReducer(
		0,
		func(v int) int { return v * v },
		func(a, v int) int { return a + v },
	)

	testTable := []struct {
		input  []int
		output int
	}{
		{input: []int{1, 2, 3, 4}, output: 30},
		{input: []int{1, 2, 3, 4}, output: 60},
		{input: []int{}, output: 60},
	}

	for i := 0; i < len(testTable); i++ {
		testCase := testTable[i]
		output := powerSum(testCase.input...)
		if testCase.output != output {
			t.Errorf("Call %d, expected output %v, actual output %v", i+1, testCase.output, output)
		}
	}
}
