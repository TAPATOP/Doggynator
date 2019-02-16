package source

import (
	"fmt"
	"testing"
)

func TestStringToResponse(t *testing.T) {
	var tests = []struct {
		input    string
		expected Response
	}{
		{
			input:    "irrelevant",
			expected: Response(DontKnowOrIrrelevant),
		},
		{
			input:    "no",
			expected: Response(No),
		},
		{
			input:    "irReleVaNt",
			expected: Response(DontKnowOrIrrelevant),
		},
		{
			input:    "yes",
			expected: Response(Yes),
		},
		{
			input:    "yea",
			expected: Response(Yes),
		},
		{
			input:    "yup",
			expected: Response(Yes),
		},
	}

	for index, test := range tests {
		t.Run(fmt.Sprintf("Response::StringToResponse(%d)", index), func(t *testing.T) {
			result := stringToResponse(test.input)
			if result != test.expected {
				createErrorWhenExpectingString(t, "StringToResponse", result.toString(), test.expected.toString())
			}
		})
	}
}
