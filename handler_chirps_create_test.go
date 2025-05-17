package main

import (
	"testing"
)

func TestChirpProfaneChecker(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No profane words",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "Single profane word lowercase",
			input:    "This is kerfuffle",
			expected: "This is ****",
		},
		{
			name:     "Single profane word uppercase",
			input:    "This is KERFUFFLE",
			expected: "This is ****",
		},
		{
			name:     "Profane word in the middle",
			input:    "foo sharbert bar",
			expected: "foo **** bar",
		},
		{
			name:     "Multiple profane words",
			input:    "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			name:     "Profane word with punctuation",
			input:    "kerfuffle, sharbert.",
			expected: "kerfuffle, sharbert.",
		},
		{
			name:     "Profane word mixed case",
			input:    "I saw a ForNaX",
			expected: "I saw a ****",
		},
		{
			name:     "Profane word as part of another word",
			input:    "sharberting is not a word",
			expected: "sharberting is not a word",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := validateChirp(tt.input)
			if result != tt.expected {
				t.Errorf("chirpProfaneChecker(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
