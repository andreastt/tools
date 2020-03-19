// +build darwin

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// this is copied from https://github.com/mariusae/tools/blob/master/f/f.go
func f(searchdir, path string, findings chan string, done chan struct{}) {
	dir, file := filepath.Split(path)
	defer close(done)

	if dir != "" {
		scan, err := mdfind(searchdir, "kind:folder", filepath.Base(dir))
		if err != nil {
			log.Fatal(err)
		}

		// filter out directories which do not match directly
		for scan.Scan() {
			match := scan.Text()
			if !strings.HasSuffix(match+"/", "/"+dir) {
				continue
			}
			if file == "" {
				findings <- match + "/"
				continue
			}

			names, err := readDirNames(match)
			if err != nil {
				log.Fatal(err)
			}
			for _, name := range names {
				matched, err := filepath.Match(file, name)
				if err != nil {
					log.Fatal(err)
				}
				if matched {
					path := filepath.Join(match, name)
					info, err := os.Stat(path)
					if err == nil && info.IsDir() {
						path += "/"
					}
					findings <- path
				}
			}
		}

		if err := scan.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		scan, err := mdfind(searchdir, fmt.Sprintf("kMDItemDisplayName == '%s'cd'", file))
		if err != nil {
			log.Fatal(err)
		}
		for scan.Scan() {
			findings <- scan.Text()
		}
		if err := scan.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func mdfind(dir string, args ...string) (*bufio.Scanner, error) {
	args = append([]string{"-onlyin", dir}, args...)
	cmd := exec.Command("mdfind", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return bufio.NewScanner(&out), nil
}

// readDirNames reads the directory named by dirname
// and returns a sorted list of directory entries.
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
