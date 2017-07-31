package main

import (
	"testing"
)

func TestDeterminePalindrome(t *testing.T) {
	tt := []struct {
		name   string
		value  string
		result bool
	}{
		{name: "single word", value: "saippuakivikauppias", result: true},
		{name: "multiple words and symbols", value: "A Man, A Plan, A Canal: Panama!", result: true},
		{name: "not a palindrome", value: "not a palindrome", result: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := &Message{Message: tc.value}
			p := m.determinePalidrome()

			if p != tc.result {
				t.Errorf("determinePalidrome() of %v should be %v; got %v", tc.value, tc.result, p)
			}
		})
	}
}
