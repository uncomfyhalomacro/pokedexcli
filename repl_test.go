package main

import "testing"

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello         world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "A walk in the park!",
			expected: []string{"a", "walk", "in", "the", "park!"},
		},
		{input: `Quite multi-lined
like   with more space!
`,
			expected: []string{"quite", "multi-lined", "like", "with", "more", "space!"},
		},
	}

	for _, testCase := range testCases {
		actual := cleanInput(testCase.input)

		for i := range actual {
			word := actual[i]
			expectedWord := testCase.expected[i]
			if expectedWord != word {
				t.Errorf("expected word: %s\ngot: %s", expectedWord, word)
			}
		}
	}
}
