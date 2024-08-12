package pkg

import (
	"fmt"
	"os"
)

// CheckErr checks for an error and prints a message with optional format and arguments before exiting the program if an error is found.
//
// It takes in an error, a format string, and optional arguments.
// If the error is not nil, it prints the formatted message and the error.
// Then it exits the program with a status code of 1.
//
// Returns:
// - error: nil if the error is nil, otherwise the error.
func CheckErr(err error, format string, a ...interface{}) error {
	if err != nil {
		if format != "" {
			msg := fmt.Sprintf(format, a...)
			fmt.Println(msg)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return nil
}
