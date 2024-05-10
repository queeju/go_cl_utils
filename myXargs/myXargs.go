/*
myXargs is a simple command-line utility similar to xargs.
It reads input from both command-line arguments and standard input,
and then executes a specified command with those arguments.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	arr := os.Args[1:]

	// Read command-line arguments
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}
	err := scanner.Err()
	CheckErr(err)

	// Execute the command with arguments
	cmd := exec.Command(arr[0], arr[1:]...)
	stdout, err := cmd.CombinedOutput()
	CheckErr(err)
	fmt.Print(string(stdout))
}

// CheckErr is a utility function to check for and handle errors.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
