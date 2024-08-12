package pkg

import (
	"fmt"
)

// CheckErr checks for an error and prints a message with optional format and arguments before exiting the program if an error is found.
//
// Parameters:
// - err: the error to check.
// - format: the format of the message to print.
// - a: the arguments to format the message.
//
// Return:
// - error: the original error if it is not nil.
func CheckErr(err error, format string, a ...interface{}) error {
	if err != nil {
		if format != "" {
			msg := fmt.Sprintf(format, a...)
			fmt.Println(msg)
		}
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}
