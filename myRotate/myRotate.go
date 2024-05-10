/*
myRotate is a utility to archive .log files into compressed tar archives.
The archive command allows archiving multiple .log files into compressed tar archives (.tar.gz). 
It takes a destination directory as input where the archived files will be stored. 
If no destination is provided, the archived files will be stored in the current directory.

Usage:
	archive -a path/to/archive/destination [file1 file2 ...]

Options:
  -a string
        path/to/archive/destination
*/
package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// main is the entry point of the program.
func main() {
	// Parse command line flags
	dest := flag.String("a", "", "path/to/archive/destination")
	flag.Parse()

	// Ensure the destination directory is valid
	destInfo, err := os.Stat(*dest)
	CheckErr(err)
	if !destInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory\n", *dest)
		return
	}

	// Get the list of files to archive
	var files []string
	if *dest == "" {
		files = os.Args[1:]
	} else {
		files = os.Args[3:]
	}

	// Check if files are provided
	if len(files) < 1 {
		panic("No files provided for archiving")
	}

	// Archive each file concurrently
	wg := new(sync.WaitGroup)
	for _, path := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			info, err := os.Stat(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			name, ok := getArchiveName(&info)
			if !ok {
				return
			}

			res, err := os.Create(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			defer res.Close()
			fmt.Println("CREATED:", name)

			err = fillArchive(&path, res, &info)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				removeArchive(&name)
				return
			}

			if *dest != "" {
				moveArchive(&name, dest)
			}
		}(path)
	}
	wg.Wait()
}

// getArchiveName generates the name for the compressed tar archive based on file information.
func getArchiveName(info *fs.FileInfo) (string, bool) {
	stamp := (*info).ModTime().Unix()
	name, found := strings.CutSuffix((*info).Name(), ".log")
	if !found {
		fmt.Fprintln(os.Stderr, "Wrong file format, only .log accepted")
		return "", false
	}
	name = fmt.Sprintf("%s_%d.tar.gz", name, stamp)
	return name, true
}

// moveArchive moves the created archive to the destination directory.
func moveArchive(name, dest *string) {
	cmd := exec.Command("mv", *name, *dest)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Print(string(stdout))
}

// removeArchive removes the archive if an error occurs during archiving.
func removeArchive(name *string) {
	cmd := exec.Command("rm", *name)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Print(string(stdout))
}

// fillArchive creates the compressed tar archive and fills it with file contents.
func fillArchive(path *string, buf io.Writer, info *fs.FileInfo) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	file, err := os.Open(*path)
	if err != nil {
		return err
	}
	defer file.Close()

	// create correct tar header
	header, err := tar.FileInfoHeader(*info, (*info).Name())
	if err != nil {
		return err
	}
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// copy file contents
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	fmt.Println("ARCHIVED:", *path)
	return nil
}

// CheckErr is a utility function to check for and handle errors.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
