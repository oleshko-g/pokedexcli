package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oleshko-g/pokedexcli/internal/pokecache"
)

const (
	locationAreaURL = "https://pokeapi.co/api/v2/location-area/"
	pokemonURL      = "https://pokeapi.co/api/v2/pokemon/"
	printConfig     = true
	printResponse   = true
)

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
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

type locationAreaByNameResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
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
var pokedex = make(map[string]Pokemon)

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
		"explore": {
			name:        "expore",
			description: "Displays pokemons which can be encountered in a location area",
			callback:    commandExpore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a named pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a caught pokemn by its name",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists caught pokemon",
			callback:    commandPokedex,
		},
	}

	s := locationAreaURL
	cfg.nextURL = &s
}

func commandUknown() error {
	return fmt.Errorf("unknown command")
}

func commandExit(...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func printLocations(l locationAreaResponse) {
	for _, v := range l.Locations {
		fmt.Printf("%v\n", v.Name)
	}
}

func printPokemons(ln locationAreaByNameResponse) {
	fmt.Println("Found Pokemon:")
	for _, v := range ln.PokemonEncounters {
		fmt.Printf("- %v\n", v.Pokemon.Name)
	}
}

func fetchData(url string) ([]byte, error) {
	if val, ok := cache.Get(url); ok {
		return val, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)

	return data, nil
}

func fetchAndPrintLocations(url string) error {
	if url == "" {
		return fmt.Errorf("URL is empty")
	}

	cfg.print()
	ld, err := fetchData(url)
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

func commandMap(...string) error {
	if cfg.nextURL == nil {
		fmt.Println("You're at the last page. Type \"mapb\" to go on the previous page")
		return nil
	}

	return fetchAndPrintLocations(*cfg.nextURL)
}

func commandExpore(s ...string) error {
	url := locationAreaURL + s[0]
	data, err := fetchData(url)
	if err != nil {
		return err
	}

	var ln locationAreaByNameResponse
	err = json.Unmarshal(data, &ln)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring \"%s\" location area...\n", ln.Name)
	printPokemons(ln)

	return nil
}

func commandMapB(...string) error {
	if cfg.prevURL == nil {
		fmt.Println("You're at the first page. Type \"map\" to go on the next page")
		return nil
	}

	return fetchAndPrintLocations(*cfg.prevURL)
}

func commandCatch(pokemonName ...string) error {
	if pokemonName[0] == "" || pokemonName == nil {
		return fmt.Errorf("pokemon name is empty")
	}
	data, err := fetchData(pokemonURL + pokemonName[0])
	if err != nil {
		return err
	}

	var p Pokemon
	err = json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)
	catchChance := rand.Int()
	if catchChance < p.BaseExperience {
		fmt.Printf("%s escaped!\n", p.Name)
		return nil
	}
	pokedex[p.Name] = p
	fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", p.Name)

	return nil
}

func printField(fieldName string, fieldValue string) {
	fmt.Printf("%s: %v\n", fieldName, fieldValue)
}

func commandInspect(pokemonName ...string) error {
	pokemon, ok := pokedex[pokemonName[0]]
	if !ok {
		fmt.Printf("You havn't caught %s yet\n", pokemonName[0])
		return nil
	}
	printField("Name", pokemon.Name)
	printField("Height", fmt.Sprintf("%d", pokemon.Height))
	printField("Weight", fmt.Sprintf("%d", pokemon.Weight))
	fmt.Println("Stats:")
	for _, v := range pokemon.Stats {
		fmt.Print("  - ")
		printField(v.Stat.Name, strconv.Itoa(v.BaseStat))
	}
	fmt.Println("Types:")
	for _, v := range pokemon.Types {
		fmt.Printf("  - %s\n", v.Type.Name)
	}
	return nil
}

func commandPokedex(...string) error {
	if len(pokedex) == 0 {
		fmt.Println("You haven't caught any pokemon yet")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for pokemon := range pokedex {
		fmt.Printf("  - %s\n", pokemon)
	}
	return nil
}

func commandHelp(...string) error {
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
	if len(words) > 1 {
		return command.callback(words[1])
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
