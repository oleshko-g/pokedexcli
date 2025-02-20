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
	fmt.Print("Pokedex > ")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		words := strings.Fields(strings.ToLower(text))
		fmt.Printf("Your command was: %s\n", words[0])
		fmt.Print("Pokedex > ")
	}
}
