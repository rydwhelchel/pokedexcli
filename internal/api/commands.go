package api

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

const (
	exitString = "Thanks for playing!"
	noNextMaps = "There are no more locations!"
	noPrevMaps = "You're at the start, there are no previous locations!"
)

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"map": {
			name:        "map",
			description: "Gets the next 20 locations",
			Callback:    CommandNextMaps,
		},
		"mapb": {
			name:        "mapb",
			description: "Gets the previous 20 locations",
			Callback:    CommandPrevMaps,
		},
		"quit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"explore": {
			name:        "explore",
			description: "Shows Pokemon you can find by exploring given area",
			Callback:    CommandExploreArea,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch the specified Pokemon",
			Callback:    CommandCatchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects the details of a specified Pokemon",
			Callback:    CommandInspectPokemon,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of Pokemon currently in your Pokedex",
			Callback:    CommandPokedex,
		},
	}
}

type CliCommand struct {
	name        string
	description string
	Callback    func(*Config, string) (bool, error)
}

type Config struct {
	next   string
	prev   string
	caught []Pokemon
}

func CommandHelp(conf *Config, arg string) (bool, error) {
	fmt.Println(buildHelpString())
	return true, nil
}

func buildHelpString() string {
	helpCmds := []string{}
	for e, v := range GetCommands() {
		helpCmds = append(helpCmds, e+"    "+v.description)
	}

	helpStr := ""
	sort.Strings(helpCmds)
	for _, v := range helpCmds {
		helpStr += v + "\n"
	}
	return helpStr
}

func CommandExit(conf *Config, arg string) (bool, error) {
	fmt.Println(exitString)
	return false, nil
}

func CommandNextMaps(conf *Config, arg string) (bool, error) {
	locs := GetNextLocations(conf)
	if len(locs) == 0 {
		fmt.Println(noNextMaps)
		return true, nil
	}
	for _, v := range locs {
		fmt.Println(v.Name)
	}
	return true, nil
}

func CommandPrevMaps(conf *Config, arg string) (bool, error) {
	locs := GetPrevLocations(conf)
	if len(locs) == 0 {
		fmt.Println(noPrevMaps)
		return true, nil
	}
	for _, v := range locs {
		fmt.Println(v.Name)
	}
	return true, nil
}

func CommandExploreArea(conf *Config, arg string) (bool, error) {
	pokemon := GetExploreArea(arg)
	if len(pokemon) > 0 {
		fmt.Println("Found Pokemon...")
		for _, p := range pokemon {
			fmt.Println("  - " + p)
		}
	} else {
		fmt.Println("No Pokemon found...")
	}
	return true, nil
}

func CommandCatchPokemon(conf *Config, pokemonName string) (bool, error) {
	if conf.caught == nil {
		conf.caught = make([]Pokemon, 0)
	}
	pokemon := GetPokemon(pokemonName)
	exp := pokemon.BaseExperience
	if roll := rand.Intn(exp); roll <= catchThresh {
		fmt.Println("You caught it!")
		fmt.Println("Adding " + pokemonName + " to your Pokedex")
		conf.caught = append(conf.caught, pokemon)
		fmt.Println("So far you have:")
		for _, poke := range conf.caught {
			fmt.Println("  - " + poke.Name)
		}
	} else {
		fmt.Println("Unlucky! It broke free!")
	}
	return true, nil
}

func CommandInspectPokemon(conf *Config, pokemonName string) (bool, error) {
	have := false
	for _, v := range conf.caught {
		if strings.ToLower(v.Name) == strings.ToLower(pokemonName) {
			printPokeFax(v)
			have = true
		}
	}
	if !have {
		fmt.Println("You haven't discovered that Pokemon!")
	}
	return true, nil
}

func printPokeFax(pokemon Pokemon) {
	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, v := range pokemon.Stats {
		fmt.Printf("  - %s: %v\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range pokemon.Types {
		fmt.Printf("  - %s\n", v.Type.Name)
	}
}

func CommandPokedex(conf *Config, arg string) (bool, error) {
	fmt.Println("Your Pokedex:")
	for _, v := range conf.caught {
		fmt.Printf("  - %s\n", v.Name)
	}
	return true, nil
}
