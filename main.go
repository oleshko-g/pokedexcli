package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	t := strings.ToLower(text)

	return strings.Fields(t)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			continue
		}
		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}
		fmt.Printf("\rYour command was: %s\n", words[0])
	}
}
