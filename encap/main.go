// Encapsulates stdin with text given by first,
// and optionally, second argument.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const EX_USAGE = 64

func encap(w io.Writer, start, s, end string) {
	fmt.Fprintf(w, "%s%s%s", start, s, end)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s <start> [<end>]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()

	var start, end string
	switch flag.NArg() {
	case 2:
		start = flag.Arg(0)
		end = flag.Arg(1)
	case 1:
		start = flag.Arg(0)
		end = flag.Arg(0)
	default:
	case 0:
		flag.Usage()
	}

	r := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	encap(os.Stdout, start, string(b), end)
}
