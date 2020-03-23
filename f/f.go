/*
f finds a file on disk.

The tool searches for the file given as the first argument
in the current working directory's file hierarchy by default.
An optional second argument lets you change the file hierarchy to traverse.

It uses the best technique available for the current system.
On macOS it uses mdfind(1) which is backed by a cache,
otherwise it falls back on the uncached find(1)
found on most standard Unix systems.
*/
package main // import "sny.no/tools/f"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const EX_USAGE = 64

func usage() {
	fmt.Fprintf(os.Stderr, "usage: f file [path]\n")
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()

	path := flag.Arg(0)
	searchdir := "."
	switch flag.NArg() {
	case 1:
	case 2:
		searchdir = flag.Arg(1)
	default:
		flag.Usage()
	}

	findings := make(chan string)
	done := make(chan struct{})
	go f(searchdir, path, findings, done)
	for {
		select {
		case file := <-findings:
			fmt.Println(relativizePath(file))
		case <-done:
			return
		}
	}
}

func relativizePath(path string) string {
	base, err := os.Getwd()
	if err != nil {
		log.Println("unable to get cwd:", err)
		return path
	}
	if strings.HasPrefix(path, base) {
		rel, err := filepath.Rel(base, path)
		if err != nil {
			log.Println("unable to get relative path:", err)
			return path
		}
		return fmt.Sprintf("./%s", rel)
	}
	return path
}
