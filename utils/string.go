package utils

import "fmt"

// divide string to max n
func DivideString(input string, n int ) string {
	if len(input) <= n {
		return input
	}

	a := []rune(input)

	return fmt.Sprintf("%s...", string(a[0:n-4]))
}
