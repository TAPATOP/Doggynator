package source

import (
	"fmt"
	"testing"
)

// TODO: Move this out of here......
func TestCustomDataBaseOfFactsConstructor(t *testing.T) {
	var tests = []struct {
		nameForMethod   string
		answers         []int
		answeredIndexes []int
		expectedCount   int
	}{
		{
			nameForMethod:   "CustomDataBaseConstructor with zeroes",
			answers:         []int{1, 0, 0, 0, 1, 2},
			answeredIndexes: []int{IsAnswered, 0, 0, IsAnswered, IsAnswered, IsAnswered},
			expectedCount:   4,
		},
		{
			nameForMethod:   "CustomDataBaseConstructor with no zeroes",
			answers:         []int{1, 0, 0, 0, 1, 2},
			answeredIndexes: []int{IsAnswered, 0, 0, 0, IsAnswered, IsAnswered},
			expectedCount:   3,
		},
		{
			nameForMethod:   "CustomDataBaseConstructor with one question 1",
			answers:         []int{0},
			answeredIndexes: []int{IsAnswered},
			expectedCount:   1,
		},
		{
			nameForMethod:   "CustomDataBaseConstructor with one question 2",
			answers:         []int{0},
			answeredIndexes: []int{0},
			expectedCount:   0,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("AskQuestion(expected:%d)", test.expectedCount), func(t *testing.T) {
			dbf := *CustomDataBaseOfFactsConstructor(test.answers, test.answeredIndexes)
			result := dbf.recordedAnswerNumber
			if result != test.expectedCount {
				createErrorWhenExpectingInt(t, test.nameForMethod, result, test.expectedCount)
			}
		})
	}
}

func TestLearn(t *testing.T) {
	var tests = []struct {
		nameForMethod string
		record        *Record
		dbf           *DataBaseOfFacts
		expected      *Record
	}{
		{
			nameForMethod: "Learn with three questions and two answers 1",
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1, 1, 1}),
					*StatisticConstructor([StatisticSize]int{2, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 5}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{0, 0, 2},
				[]int{0, 0, IsAnswered},
			),
			expected: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1, 1, 1}),
					*StatisticConstructor([StatisticSize]int{2, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 6}),
				}),
		},
		{
			nameForMethod: "Learn with three questions and two answers 2",
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1, 1, 1}),
					*StatisticConstructor([StatisticSize]int{2, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 5}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{1, 0, 0},
				[]int{IsAnswered, IsAnswered, 0},
			),
			expected: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1, 2, 1}),
					*StatisticConstructor([StatisticSize]int{3, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 5}),
				}),
		},
		{
			nameForMethod: "Learn with three questions and two answers and reduce",
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1, DataMaximumBeforeReduce, 1}),
					*StatisticConstructor([StatisticSize]int{2, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 5}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{1, 0, 0},
				[]int{IsAnswered, IsAnswered, 0},
			),
			expected: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{1 / ReductionFactor, DataMaximumBeforeReduce / ReductionFactor, 1 / ReductionFactor}),
					*StatisticConstructor([StatisticSize]int{3, 5, 0}),
					*StatisticConstructor([StatisticSize]int{4, 5, 5}),
				}),
		},
		{
			nameForMethod: "Learn with three questions, two answers, too many answers and reduce",
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{DataMaximumBeforeReduce / 3, DataMaximumBeforeReduce / 2, DataMaximumBeforeReduce / 5}),
					*StatisticConstructor([StatisticSize]int{0, DataMaximumBeforeReduce - 1, 0}),
					*StatisticConstructor([StatisticSize]int{12, DataMaximumBeforeReduce, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{1, 0, 0},
				[]int{IsAnswered, 0, IsAnswered},
			),
			expected: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{DataMaximumBeforeReduce / (3 * ReductionFactor), DataMaximumBeforeReduce / (2 * ReductionFactor), DataMaximumBeforeReduce / (5 * ReductionFactor)}),
					*StatisticConstructor([StatisticSize]int{0, DataMaximumBeforeReduce - 1, 0}),
					*StatisticConstructor([StatisticSize]int{12 / ReductionFactor, DataMaximumBeforeReduce / ReductionFactor, 2 / ReductionFactor}),
				}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.nameForMethod), func(t *testing.T) {
			lb := LearningMechanismConstructor(test.dbf)
			lb.learn(test.record)
			if test.record.ToString() != test.expected.ToString() {
				createErrorWhenExpectingString(t, test.nameForMethod, test.record.ToString(), test.expected.ToString())
			}
		})
	}
}

// TODO: Test reduce()
