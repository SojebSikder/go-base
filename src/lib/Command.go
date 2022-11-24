package lib

import (
	"fmt"
	"strings"
)

// YesNoPrompt asks yes/no questions using the label.
func YesNoPrompt(label string) bool {

	// show the lebel
	fmt.Print(label)
	// take input
	var input string
	fmt.Scan(&input)

	input = strings.ToLower(input)
	if input == "y" || input == "yes" {
		return true

	} else if input == "n" || input == "no" {
		return false
	}

	return YesNoPrompt(label)
}
