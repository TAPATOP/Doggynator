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
	ie := InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(0, Response(Yes))
	exp := 0
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q0 Yes", n, exp)
	}

	ie = InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(0, Response(No))
	exp = 2
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q0 No", n, exp)
	}

	ie = InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(0, Response(DontKnowOrIrrelevant))
	exp = 1
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q0 DK", n, exp)
	}

	ie = InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(1, Response(DontKnowOrIrrelevant))
	exp = 2
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q1 Yes", n, exp)
	}

	ie = InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(1, Response(No))
	exp = 0
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q1 No", n, exp)
	}

	ie = InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{})
	ie.processResponse(1, Response(DontKnowOrIrrelevant))
	exp = 2
	if n := ie.getBestGuessIndex(); n != exp {
		createError(t, "Get Best Guess Index Q1 DK", n, exp)
	}
}

func createError(t *testing.T, nameOfMethod string, returned, expected int) {
	t.Error(nameOfMethod + " returned " + strconv.Itoa(returned) + " instead of " + strconv.Itoa(expected))
}
