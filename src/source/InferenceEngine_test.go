package source

import (
	"fmt"
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
		{nameForMethod: "Get Best Guess Index Q0 Yes", input: 0, response: Yes, expected: 0},
		{nameForMethod: "Get Best Guess Index Q0 No", input: 0, response: No, expected: 2},
		{nameForMethod: "Get Best Guess Index Q0 DK", input: 0, response: DontKnowOrIrrelevant, expected: 1},
		{nameForMethod: "Get Best Guess Index Q1 Yes", input: 1, response: Yes, expected: 0},
		{nameForMethod: "Get Best Guess Index Q1 No", input: 1, response: No, expected: 0},
		{nameForMethod: "Get Best Guess Index Q1 DK", input: 1, response: DontKnowOrIrrelevant, expected: 2},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("TestProcessResponse(%d)(%s)", test.input, test.response.toString()), func(t *testing.T) {
			ie := InferenceEngineConstructor(records, []string{}, &DataBaseOfFacts{}, &DefaultRandomGenerator{})
			ie.processResponse(test.input, test.response)
			if result := ie.getBestGuessIndex(); result != test.expected {
				createError(t, test.nameForMethod, result, test.expected)
			}
		})
	}
}

func TestAskQuestion(t *testing.T) {
	records := []Record{
		*RecordConstructor(
			"one",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{1, 1, 0}),
				*StatisticConstructor([StatisticSize]int{100, 5, 0}),
				*StatisticConstructor([StatisticSize]int{35, 30, 0}),
			}),
		*RecordConstructor(
			"two",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{1, 1, 1}),
				*StatisticConstructor([StatisticSize]int{3, 2, 0}),
				*StatisticConstructor([StatisticSize]int{3, 7, 0}),
			}),
		*RecordConstructor(
			"three",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{1, 1, 1}),
				*StatisticConstructor([StatisticSize]int{100, 5, 0}),
				*StatisticConstructor([StatisticSize]int{4, 5, 35}),
			}),
	}

	var tests = []struct {
		nameForMethod string
		expected      int
	}{
		{nameForMethod: "AskQuestion with best entropy", expected: 2},
		{nameForMethod: "AskQuestion with 2nd best entropy", expected: 0},
		{nameForMethod: "AskQuestion with 3rd best entropy", expected: 1},
		{nameForMethod: "AskQuestion with all questions already asked", expected: -1},
	}

	questions := []string{"q1", "q2", "q3"}

	ie := InferenceEngineConstructor(records, questions, DataBaseOfFactsConstructor(len(questions)), &FakeRandomGenerator{})

	for _, test := range tests {
		t.Run(fmt.Sprintf("AskQuestion(expected:%d)", test.expected), func(t *testing.T) {
			result := ie.askQuestion()
			if result != test.expected {
				createError(t, test.nameForMethod, result, test.expected)
			}
			if result != -1 {
				ie.dbf.record(0, result)
			}
		})
	}
}

func TestConcludeAnAnswer(t *testing.T) {
	records := []Record{
		*RecordConstructor(
			"one",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{13, 0, 0}),
				*StatisticConstructor([StatisticSize]int{3, 35, 0}),
				*StatisticConstructor([StatisticSize]int{35, 35, 2000}),
				*StatisticConstructor([StatisticSize]int{35, 35, 0}),
				*StatisticConstructor([StatisticSize]int{35, 35, 0}),
			}),
		*RecordConstructor(
			"two",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{13, 1, 1}),
				*StatisticConstructor([StatisticSize]int{2, 0, 0}),
				*StatisticConstructor([StatisticSize]int{2, 2, 0}),
				*StatisticConstructor([StatisticSize]int{2, 2, 0}),
				*StatisticConstructor([StatisticSize]int{2, 2, 0}),
			}),
		*RecordConstructor(
			"three",
			[]Statistic{
				*StatisticConstructor([StatisticSize]int{1, 3, 55}),
				*StatisticConstructor([StatisticSize]int{4, 1, 35}),
				*StatisticConstructor([StatisticSize]int{4, 1, 1}),
				*StatisticConstructor([StatisticSize]int{2, 1, 35}),
				*StatisticConstructor([StatisticSize]int{1, 2, 35}),
			}),
	}

	type Input struct {
		index    int
		response Response
	}

	var tests = []struct {
		nameForMethod            string
		input                    []Input
		enquiriesSinceLastAnswer int
		expected                 int
	}{
		{
			nameForMethod: "Conclude the second answer",
			input: []Input{
				{0, DontKnowOrIrrelevant},
				{1, DontKnowOrIrrelevant},
				{2, No},
				{3, No},
				{4, No},
			},
			enquiriesSinceLastAnswer: MinimumAnsweredQuestions + MaximumIntervalBetweenAnswers,
			expected:                 1,
		},
		{
			nameForMethod: "Conclude the first answer",
			input: []Input{
				{0, Yes},
				{1, No},
				{2, DontKnowOrIrrelevant},
				{3, Yes},
				{4, Yes},
			},
			enquiriesSinceLastAnswer: MinimumAnsweredQuestions + MaximumIntervalBetweenAnswers + 2,
			expected:                 0,
		},
		{
			nameForMethod: "Conclude the third highly probable answer",
			input: []Input{
				{0, Yes},
				{1, No},
				{2, Yes},
				{3, DontKnowOrIrrelevant},
				{4, DontKnowOrIrrelevant},
			},
			enquiriesSinceLastAnswer: MaximumIntervalBetweenAnswers - 1,
			expected:                 2,
		},
	}

	for _, test := range tests {
		ie := InferenceEngineConstructor(records, []string{}, DataBaseOfFactsConstructor(len(records[0].statistics)), &FakeRandomGenerator{})
		for _, input := range test.input {
			ie.processResponse(input.index, input.response)
		}
		ie.enquiriesSinceLastAnswer = test.enquiriesSinceLastAnswer
		_, resultIndex := ie.concludeAnAnswer()
		if resultIndex != test.expected {
			createError(t, test.nameForMethod, resultIndex, test.expected)
		}
	}
}

// Helper stuff //

func createError(t *testing.T, nameOfMethod string, returned, expected int) {
	t.Error(nameOfMethod + " returned " + strconv.Itoa(returned) + " instead of " + strconv.Itoa(expected))
}

type FakeRandomGenerator struct{}

func (obj *FakeRandomGenerator) Intn(limit int) int {
	return 0
}
