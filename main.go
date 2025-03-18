package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oleshko-g/pokedexcli/internal/pokecache"
)

const (
	locationAreaURL = "https://pokeapi.co/api/v2/location-area/"
	printConfig     = true
	printResponse   = true
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	nextURL *string
	prevURL *string
}

func (c config) print() {
	if !printConfig {
		return
	}

	fmt.Println("")
	fmt.Println("## Config")

	if cfg.prevURL != nil {
		fmt.Printf("prevURL: %#v\n", *(cfg.prevURL))
	} else {
		fmt.Printf("prevURL: %#v\n", cfg.prevURL)
	}

	if cfg.nextURL != nil {
		fmt.Printf("nextURL: %#v\n", *(cfg.nextURL))
	} else {
		fmt.Printf("nextURL: %#v\n", cfg.nextURL)
	}

	fmt.Println("")
}

type locationAreaResponse struct {
	Locations []struct {
		Name string `json:"name"`
	} `json:"results"`
	PreviousURL *string `json:"previous"`
	NextURL     *string `json:"next"`
}

func (l locationAreaResponse) print() {
	if !printResponse {
		return
	}

	fmt.Println("")
	fmt.Println("## Response")

	if l.PreviousURL != nil {
		fmt.Printf("PreviousURL: %#v\n", *(l.PreviousURL))
	} else {
		fmt.Printf("PreviousURL: %#v\n", l.PreviousURL)
	}

	if l.NextURL != nil {
		fmt.Printf("NextURL: %#v\n", *(l.NextURL))
	} else {
		fmt.Printf("NextURL: %#v\n", l.NextURL)
	}

	fmt.Println("")
}

var commandRegistry map[string]cliCommand
var cfg config
var cache = pokecache.NewCache(5 * time.Second)

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
		"map": {
			name:        "map",
			description: "Displays names of 20 location next areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays names of 20 location previous areas in the Pokemon world",
			callback:    commandMapB,
		},
	}

	s := locationAreaURL
	cfg.nextURL = &s
}

func commandUknown() error {
	return fmt.Errorf("unknown command")
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func printLocations(l locationAreaResponse) {
	for _, v := range l.Locations {
		fmt.Printf("%v\n", v.Name)
	}
}

func fetchLocationsData(url string) ([]byte, error) {
	if val, ok := cache.Get(url); ok {
		return val, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	locationsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, locationsData)

	return locationsData, nil
}

func fetchAndPrintLocations(url string) error {
	if url == "" {
		return fmt.Errorf("URL is empty")
	}

	cfg.print()
	ld, err := fetchLocationsData(url)
	if err != nil {
		return err
	}

	var l locationAreaResponse
	if err := json.Unmarshal(ld, &l); err != nil {
		return err
	}
	l.print()

	printLocations(l)

	cfg.prevURL = l.PreviousURL
	cfg.nextURL = l.NextURL

	cfg.print()

	return nil
}

func commandMap() error {
	if cfg.nextURL == nil {
		fmt.Println("You're at the last page. Type \"mapb\" to go on the previous page")
		return nil
	}

	return fetchAndPrintLocations(*cfg.nextURL)
}

func commandMapB() error {
	if cfg.prevURL == nil {
		fmt.Println("You're at the first page. Type \"map\" to go on the next page")
		return nil
	}

	return fetchAndPrintLocations(*cfg.prevURL)
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

		err := processCommand(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func main() {
	repl(bufio.NewScanner(os.Stdin))
}
