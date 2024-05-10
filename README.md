# Custom Command-line Utilities

This repository contains a collection of custom command-line utilities implemented in Go. Each utility serves a specific purpose and provides functionalities similar to popular Unix-like commands.
This project was developed as a part of School 21 curriculum.

To use each of the utilities, navigate to the appropriate folder and run 

`go build`

## Utilities

### **1. myFind**

The `myFind` utility locates files, directories, and symbolic links within a specified directory, supporting recursive searching and various filtering options. Users can customize the output to focus on specific file types or extensions. The program resolves symbolic links, indicates broken links as `[broken]` in the output, and skips files and directories inaccessible to the current user.

#### Usage:

`./myFind [options] /path/to/directory` 

#### Options:

-   `-f`: Print only regular files.
-   `-d`: Print only directories.
-   `-sl`: Print only symbolic links.
-   `-ext`: Specify file extension to filter results (works only with `-f` option).

### **2. myRotate**

The `myRotate` utility is used to archive `.log` files into compressed tar archives. It allows archiving multiple `.log` files into compressed tar archives (`.tar.gz`). The destination directory for the archived files can be specified. If no destination is provided, the archived files will be stored in the current directory.

#### Usage:

`./myRotate [-a path/to/archive/destination] file1 [file2 ...]` 

#### Options:

-   `-a string`: Path to the archive destination.

### **3. myWc**

The `myWc` utility is similar to the `wc` command in Unix-like operating systems. It counts lines, words, and characters in text files.

#### Usage:

`./myWc [options] file1 [file2 ...]` 

#### Options:

-   `-l`: Count lines
-   `-w`: Count words
-   `-m`: Count characters

If no options are specified, words are counted by default.

### **4. myXargs**

The `myXargs` utility is a simple command-line utility similar to `xargs`. It reads input from both command-line arguments and standard input, and then executes a specified command with those arguments.
