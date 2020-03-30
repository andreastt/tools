/*
Harmonises clipboards across systems and networks.

Text can be piped onto the primary clipboard:

	% echo "foo" | clip
	% clip
	foo

And from file:

	% clip <bar.txt
	% clip
	bar

Or from plain arguments:

	% clip baz
	% clip
	baz

Remaining work to be done:

- support for X windowing system/Linux
- support for clipboards across ssh
*/
package main // import "sny.no/tools/clip"

import (
	"log"
	"fmt"
	"strings"
	"flag"
	"io/ioutil"
	"os"
)

const EX_USAGE = 64

func usage() {
	fmt.Fprintf(os.Stderr, "usage: clip [<stdin>|args...]\n")
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		copy(os.Stdin)
	} else if len(os.Args) > 1 {
		in := strings.Join(os.Args[1:], " ")
		copy(strings.NewReader(in))
	} else {
		r, err := paste()
		if err != nil {
			log.Fatal(err)
		}
		ba, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		s := strings.TrimSpace(string(ba))
		fmt.Println(s)
	}
}
