package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rydwhelchel/pokedexcli/internal/api"
)

const (
	badInput = "Bad input, please try again"
)

var conf api.Config

func main() {
	commands := api.GetCommands()
	scanner := bufio.NewScanner(os.Stdin)
	conf = api.Config{}
	eventLoop(commands, scanner)
}

func eventLoop(commands map[string]api.CliCommand, scanner *bufio.Scanner) {
	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")
		newInput := scanner.Scan()
		if newInput {
			input := scanner.Text()
			splitput := strings.Fields(input)
			lengthSplit := len(splitput)
			if lengthSplit == 0 {
				continue
			}
			if lengthSplit == 1 {
				splitput = append(splitput, "")
			}
			if lengthSplit > 2 {
				fmt.Println("Not sure what to do with that input...")
				continue
			}

			if command, ok := commands[splitput[0]]; ok {
				cont, err := command.Callback(&conf, splitput[1])
				if err != nil {
					fmt.Println(err)
					return
				}
				if !cont {
					return
				}
			} else {
				fmt.Println(badInput)
			}
		}
	}
}
