package source

import (
	"fmt"
	"testing"
)

func TestRawStatisticConstructor(t *testing.T) {
	var tests = []struct {
		nameForMethod string
		input         string
		expected      *Statistic
	}{
		{
			nameForMethod: "Construct Statistic with three single digit numbers",
			input:         "0 2 2",
			expected:      StatisticConstructor([StatisticSize]int{0, 2, 2}),
		},
		{
			nameForMethod: "Construct Statistic with three multi digit numbers",
			input:         "25 1 2209",
			expected:      StatisticConstructor([StatisticSize]int{25, 1, 2209}),
		},
		{
			nameForMethod: "Construct Statistic with three multi digit numbers, one of which is weird",
			input:         "25 01 2209",
			expected:      StatisticConstructor([StatisticSize]int{25, 1, 2209}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s(expected:%d)", test.nameForMethod, test.expected), func(t *testing.T) {
			result, err := RawStatisticConstructor(test.input)
			if err != nil || result.ToString() != test.expected.ToString() {
				createErrorWhenExpectingString(t, test.nameForMethod, result.ToString(), test.expected.ToString())
			}
		})
	}
}

func purify(statistic *Statistic, _ error) *Statistic {
	return statistic
}

func TestMostProbableAnswer(t *testing.T) {
	var tests = []struct {
		nameForMethod string
		input         *Statistic
		expected      Response
	}{
		{
			nameForMethod: "Test last attribute when difference between best and worst is 1",
			input:         purify(RawStatisticConstructor("2 2 3")),
			expected:      Response(DontKnowOrIrrelevant),
		},
		{
			nameForMethod: "Test first attribute with bigger differences",
			input:         purify(RawStatisticConstructor("52 2 12")),
			expected:      Response(Yes),
		},
		{
			nameForMethod: "Test second attribute with bigger differences",
			input:         purify(RawStatisticConstructor("5 22 12")),
			expected:      Response(No),
		},
		{
			nameForMethod: "Test third attribute with bigger differences",
			input:         purify(RawStatisticConstructor("5 2 12")),
			expected:      Response(DontKnowOrIrrelevant),
		},
		{
			nameForMethod: "Test with equal attributes",
			input:         purify(RawStatisticConstructor("1 1 1")),
			expected:      Response(Yes),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s(expected:%d)", test.nameForMethod, test.expected), func(t *testing.T) {
			result := test.input.mostProbableAnswer()
			if result != test.expected {
				createErrorWhenExpectingString(t, test.nameForMethod, result.ToString(), test.expected.ToString())
			}
		})
	}
}
