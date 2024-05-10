/*
myFind utility locates files, directories, and symbolic links within a pecified directory,
supporting recursive searching and various filtering options.
Users can customize the output to focus on specific file types or extensions.
The program resolves symbolic links, indicates broken links as [broken]
in the output, and skips files and directories inaccessible to the current user.

Usage:
	myFind [options] /path/to/directory

Options:
	-f: Print only regular files.
	-d: Print only directories.
	-sl: Print only symbolic links.
	-ext: Specify file extension to filter results (works only with -f option).
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// main function is the entry point of the myFind utility.
func main() {
	var exts extentions
	slMode := flag.Bool("sl", false, "Print symlinks")
	fMode := flag.Bool("f", false, "Print files")
	dirMode := flag.Bool("d", false, "Print directories")
	flag.Var(&exts, "ext", "Specify file extension to filter results")
	flag.Parse()

	// Validate command-line flags.
	if !*fMode && len(exts) > 0 {
		CheckErr(errors.New("Need -f to specify extensions"))
	}

	// Set default flags if none are provided.
	if !*slMode && !*fMode && !*dirMode {
		*slMode, *fMode, *dirMode = true, true, true
	}

	// Retrieve the root directory from command-line arguments.
	var root string
	if len(flag.Args()) != 1 {
		panic("One directory for searching must be provided")
	} else {
		root = flag.Arg(0)
	}

	// Ensure the root directory is valid
	rootInfo, err := os.Stat(root)
	CheckErr(err)
	if !rootInfo.IsDir() {
		panic(fmt.Sprintf("%s is not a directory", root))
	}

	// Walk through the directory and print entities based on the provided options.
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.Mode()&fs.ModePerm&fs.FileMode(0400) != 0 {
			printEntity(&path, info, dirMode, fMode, slMode, &exts)
		} else {
			// Print permission denied message for inaccessible files or directories.
			fmt.Fprintf(os.Stderr, "%s: Permission denied\n", path)
		}
		return nil
	})
	CheckErr(err)
}

// CheckErr is a utility function to check for and handle errors.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// extentions is a custom type for handling file extensions.
type extentions []string

// String method returns a string representation of the extensions.
func (e *extentions) String() string {
	return fmt.Sprintf("%s", *e)
}

// Set method sets the value of the extensions based on the provided input.
func (e *extentions) Set(val string) error {
	if !validExt(val) {
		return fmt.Errorf("Invalid extension")
	}
	*e = append(*e, fmt.Sprintf(".%s", val))
	return nil
}

// validExt is a utility function to validate file extensions using regular expressions.
var validExt = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// printEntity prints the entity (file, directory, or symbolic link) based on the specified options.
func printEntity(path *string, info os.FileInfo, dirMode, fMode, slMode *bool, exts *extentions) {
	if info.IsDir() && *dirMode {
		fmt.Println(*path)
	} else if info.Mode().IsRegular() && *fMode {
		if len(*exts) > 0 {
			processExt(path, exts)
		} else {
			fmt.Println(*path)
		}
	} else if info.Mode()&fs.ModeSymlink != 0 && *slMode {
		origFile, err := filepath.EvalSymlinks(info.Name())
		if err != nil {
			origFile = "[broken]"
		}
		fmt.Printf("%s -> %s\n", *path, origFile)
	}
}

// processExt prints the entity if it matches any of the specified file extensions.
func processExt(path *string, exts *extentions) {
	for _, ext := range *exts {
		if strings.HasSuffix(*path, ext) {
			fmt.Println(*path)
		}
	}
}
