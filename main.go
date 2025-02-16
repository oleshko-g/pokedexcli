package main

import (
	"strings"
)

func cleanInput(text string) []string {
	t := strings.ToLower(text)

	return strings.Fields(t)
}

func main() {
	print("Hello, World!")
}
