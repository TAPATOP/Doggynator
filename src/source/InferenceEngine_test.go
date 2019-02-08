package source

import (
	"strconv"
	"testing"
)

func TestProcessResponse(t *testing.T) {
	records := []Record{
		*RecordConstructor(
			"one",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{1, 0, 0}),
				*StatisticConstructor([StatisticSize]int{35, 35, 0}),
			}),
		*RecordConstructor(
			"two",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{13, 1, 1}),
				*StatisticConstructor([StatisticSize]int{2, 2, 0}),
			}),
		*RecordConstructor(
			"three",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{0, 3, 0}),
				*StatisticConstructor([StatisticSize]int{4, 1, 35}),
			}),
	}
	var tests = []struct {
		nameForMethod string
		input         int
		response      Response
		expected      int
	}{
		{nameForMethod: "Get Best Guess Index Q0 Yes", input: 0, response: (Yes), expected: 0},
		{nameForMethod: "Get Best Guess Index Q0 No", input: 0, response: (No), expected: 2},
		{nameForMethod: "Get Best Guess Index Q0 DK", input: 0, response: (DontKnowOrIrrelevant), expected: 1},
		{nameForMethod: "Get Best Guess Index Q1 Yes", input: 1, response: (Yes), expected: 0},
		{nameForMethod: "Get Best Guess Index Q1 No", input: 1, response: (No), expected: 0},
		{nameForMethod: "Get Best Guess Index Q1 DK", input: 1, response: (DontKnowOrIrrelevant), expected: 2},
	}
	for _, test := range tests {
		ie := InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
		ie.processResponse(test.input, test.response)
		if result := ie.getBestGuessIndex(); result != test.expected {
			createError(t, test.nameForMethod, result, test.expected)
		}
	}
}

//
//func TestGetHighestEntropyIndex(t *testing.T) {
//	records := []Record{
//		*RecordConstructor(
//			"one",
//			[]Statistic{
//				*StatisticConstructor([StatisticSize]int{1, 0, 0}),
//				*StatisticConstructor([StatisticSize]int{35, 35, 0}),
//			}),
//		*RecordConstructor(
//			"two",
//			[]Statistic{
//				*StatisticConstructor([StatisticSize]int{13, 1, 1}),
//				*StatisticConstructor([StatisticSize]int{2, 2, 0}),
//			}),
//		*RecordConstructor(
//			"three",
//			[]Statistic{
//				*StatisticConstructor([StatisticSize]int{0, 3, 0}),
//				*StatisticConstructor([StatisticSize]int{4, 1, 35}),
//			}),
//	}
//
//	ie := InferenceEngineConstructor(records, )
//
//}

func createError(t *testing.T, nameOfMethod string, returned, expected int) {
	t.Error(nameOfMethod + " returned " + strconv.Itoa(returned) + " instead of " + strconv.Itoa(expected))
}
