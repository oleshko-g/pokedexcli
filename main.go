package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandRegistry map[string]cliCommand

func init() {
	commandRegistry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func commandUknown() error {
	fmt.Println("Unknown command")
	return nil
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range commandRegistry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func cleanInput(text string) []string {
	t := strings.ToLower(text)

	return strings.Fields(t)
}

func processCommand(s string) error {
	words := cleanInput(s)
	if len(words) == 0 {
		return nil
	}
	command, ok := commandRegistry[words[0]]
	if !ok {
		return commandUknown()
	}
	return command.callback()
}

func repl(r *bufio.Scanner) {
	for {
		fmt.Print("Pokedex > ")
		r.Scan()
		text := r.Text()
		if text == "" {
			continue
		}
		processCommand(text)
	}
}

func main() {
	repl(bufio.NewScanner(os.Stdin))
}
