package source

import (
	"fmt"
	"testing"
)

func TestExplain(t *testing.T) {
	var tests = []struct {
		nameForMethod       string
		questions           []string
		record              *Record
		dbf                 *DataBaseOfFacts
		expectedExplanation string
		expectedSurprise    string
	}{
		{
			nameForMethod: "Explanation of 2 out of 3 correctly answered questions",
			questions:     []string{"q1", "q2", "q3"},
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{2, 2, 2}),
					*StatisticConstructor([StatisticSize]int{0, 1, 0}),
					*StatisticConstructor([StatisticSize]int{1, 0, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{0, 1, 2},
				[]int{0, IsAnswered, IsAnswered},
			),
			expectedExplanation: generateExplanationLn(Response(No), "q2") +
				generateExplanation(Response(DontKnowOrIrrelevant), "q3"),
			expectedSurprise: "",
		},
		{
			nameForMethod: "Explanation of 2 out of 3 answered questions, one wrong",
			questions:     []string{"q1", "q2", "q3"},
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{2, 2, 2}),
					*StatisticConstructor([StatisticSize]int{0, 1, 0}),
					*StatisticConstructor([StatisticSize]int{1, 0, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{0, 1, 1},
				[]int{0, IsAnswered, IsAnswered},
			),
			expectedExplanation: generateExplanation(Response(No), "q2"),
			expectedSurprise:    generateSurprised(Response(No), Response(DontKnowOrIrrelevant), "q3"),
		},
		{
			nameForMethod: "Explanation of 3 out of 3 answered questions, all wrong",
			questions:     []string{"q1", "q2", "q3"},
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{3, 2, 2}),
					*StatisticConstructor([StatisticSize]int{0, 1, 0}),
					*StatisticConstructor([StatisticSize]int{1, 0, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{1, 0, 0},
				[]int{IsAnswered, IsAnswered, IsAnswered},
			),
			expectedExplanation: "",
			expectedSurprise: generateSurprisedLn(Response(No), Response(Yes), "q1") +
				generateSurprisedLn(Response(Yes), Response(No), "q2") +
				generateSurprised(Response(Yes), Response(DontKnowOrIrrelevant), "q3"),
		},
		{
			nameForMethod: "Explanation of 3 out of 3 answered questions, all correct",
			questions:     []string{"q1", "q2", "q3"},
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{3, 2, 2}),
					*StatisticConstructor([StatisticSize]int{0, 1, 3}),
					*StatisticConstructor([StatisticSize]int{1, 0, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{0, 2, 2},
				[]int{IsAnswered, IsAnswered, IsAnswered},
			),
			expectedExplanation: generateExplanationLn(Response(Yes), "q1") +
				generateExplanationLn(Response(DontKnowOrIrrelevant), "q2") +
				generateExplanation(Response(DontKnowOrIrrelevant), "q3"),
			expectedSurprise: "",
		},
		{
			nameForMethod: "Explanation of 0 out of 3 answered questions",
			questions:     []string{"q1", "q2", "q3"},
			record: RecordConstructor(
				"three",
				[]Statistic{
					*StatisticConstructor([StatisticSize]int{3, 2, 2}),
					*StatisticConstructor([StatisticSize]int{0, 1, 0}),
					*StatisticConstructor([StatisticSize]int{1, 0, 2}),
				}),
			dbf: CustomDataBaseOfFactsConstructor(
				[]int{1, 0, 0},
				[]int{0, 0, 0},
			),
			expectedExplanation: "",
			expectedSurprise:    "",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.nameForMethod), func(t *testing.T) {
			em := ExplainingMechanismConstructor(test.questions, test.dbf)
			explanation, surprised := em.explain(test.record)
			if *explanation != test.expectedExplanation {
				createErrorWhenExpectingString(t, test.nameForMethod, *explanation, test.expectedExplanation)
			}
			if *surprised != test.expectedSurprise {
				createErrorWhenExpectingString(t, test.nameForMethod, *surprised, test.expectedSurprise)
			}
		})
	}
}
