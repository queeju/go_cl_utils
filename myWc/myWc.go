/*
myWc is a utility similar to the `wc` command in Unix-like operating systems.
It counts lines, words, and characters in text files.

Usage:
	myWc [options] file1 [file2 ...]

Options:
	-l    Count lines
	-w    Count words
	-m    Count characters

If no options are specified, words are counted by default.*/

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

var mutex sync.Mutex

// main is the entry point of the program.
func main() {
	// Command line flag declarations
	lMode := flag.Bool("l", false, "Count lines")
	wMode := flag.Bool("w", false, "Count words")
	mMode := flag.Bool("m", false, "Count characters")
	flag.Parse()

	// Validate flags and determine the mode
	mode, err := validateFlags(lMode, wMode, mMode)
	CheckErr(err)

	// Extract filenames from command line arguments
	if len(flag.Args()) < 1 {
		panic("No files provided")
	}
	files := make(map[string]int)
	for _, el := range flag.Args() {
		files[el] = 0
	}

	// Process files concurrently
	wg := new(sync.WaitGroup)
	for path := range files {
		wg.Add(1)
		go processFile(path, mode, &files, wg)
	}
	wg.Wait()

	// Print results
	for key, val := range files {
		fmt.Printf("%d\t%s\n", val, key)
	}
}

// handleFileErr outputs error message into Stderr and deletes invalid path from the map
func handleFileErr(path *string, files *map[string]int, err error) {
	fmt.Fprintln(os.Stderr, err)
	delete(*files, *path)
}

// processFile reads the file at the given path and counts lines, words, or characters based on the specified mode.
func processFile(path string, mode int8, files *map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	var res int
	var err error

	// Ensure the file path is valid
	info, err := os.Stat(path)
	if err != nil {
		handleFileErr(&path, files, err)
		return
	}

	if !info.Mode().IsRegular() {
		fmt.Fprintf(os.Stderr, "%s is not a file\n", path)
		delete(*files, path)
		return
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		handleFileErr(&path, files, err)
		return
	}
	defer file.Close()

	// Read the file line by line and count based on the mode
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		switch mode {
		case 3: // lines
			res++
		case 2: // words
			res += len(strings.Fields(scanner.Text()))
		case 1: // characters
			res += utf8.RuneCountInString(scanner.Text()) + 1
		}
	}

	// Update the map with the result
	mutex.Lock()
	(*files)[path] = res
	mutex.Unlock()
}

// CheckErr is a utility function to panic if an error is not nil.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// validateFlags validates the command line flags and determines the mode of operation.
func validateFlags(lMode, wMode, mMode *bool) (int8, error) {
	// Count the number of flags that are true
	var count, mode int8
	if *lMode {
		count++
		mode = 3
	}
	if *wMode {
		count++
		mode = 2
	}
	if *mMode {
		count++
		mode = 1
	}

	// Check if more than one flag is true
	if count > 1 {
		return 0, errors.New("Only one of -l, -w, -m can be specified")
	} else if count == 0 {
		mode = 2
	}
	return mode, nil
}
