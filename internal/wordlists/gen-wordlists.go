//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	name := os.Args[1]
	src := os.Args[2]
	dest := os.Args[3]

	f, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening source file %s: %s\n", src, err)
		os.Exit(2)
	}
	defer f.Close()

	d, err := os.Create(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening destination file %s: %s\n", dest, err)
		os.Exit(2)
	}
	defer d.Close()

	fmt.Fprint(d, "package wordlists\n\n")
	fmt.Fprintf(d, "import ()\n\n")
	fmt.Fprintf(d, "var %sWords = []string{\n", name)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		fmt.Fprintf(d, "\t\"%s\",\n", l)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin: ", err)
		os.Exit(1)
	}

	fmt.Fprintln(d, "}")
}
