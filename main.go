package main
import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for ;; {
		fmt.Print("Pokedex > ")

		var userInput []string
		if scanner.Scan() {
			userInput = cleanInput(scanner.Text())
		}
		fmt.Printf("Your command was: %s\n", userInput[0])
	}
}
