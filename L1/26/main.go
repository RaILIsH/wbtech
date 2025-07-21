package main

import (
	"fmt"
	"strings"
)

func hasAllUniqueChars(s string) bool {
	seen := make(map[rune]bool)
	lowerStr := strings.ToLower(s)

	for _, char := range lowerStr {
		if seen[char] {
			return false
		}
		seen[char] = true
	}
	return true
}

func main() {
	testCases := []string{
		"abcd",
		"abCdefAaf",
		"aabcd",
	}

	for _, testCase := range testCases {
		fmt.Printf("%q -> %t\n", testCase, hasAllUniqueChars(testCase))
	}
}
