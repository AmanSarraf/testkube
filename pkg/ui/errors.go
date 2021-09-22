package ui

import (
	"fmt"
	"os"

	"github.com/bclicn/color"
)

func ExitOnError(item string, errors ...error) {
	printAndExit(item, true, errors...)
}

func PrintOnError(item string, errors ...error) {
	printAndExit(item, false, errors...)
}

func printAndExit(item string, exitOnError bool, errors ...error) {
	if len(errors) > 0 && hasErrors(errors...) {
		for _, err := range errors {
			if err != nil {
				fmt.Printf("%s %s (error: %s)\n\n", color.LightRed("⨯"), color.Red(item), err)
				if exitOnError {
					os.Exit(1)
				}
			}
		}
	}

	if Verbose {
		fmt.Printf("%s %s\n", color.Blue("\xE2\x9C\x94"), color.Green(item))
	}
}

func hasErrors(errors ...error) bool {
	if len(errors) > 0 {
		for _, err := range errors {
			if err != nil {
				return true
			}
		}
	}

	return false
}
