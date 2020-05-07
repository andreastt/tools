// +build linux

package main

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
)

func f(searchdir, file string, findings chan string, done chan struct{}) {
	scan, err := find(searchdir, file)
	if err != nil {
		log.Fatal(err)
	}
	for scan.Scan() {
		findings <- scan.Text()
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
	close(done)
}

func find(dir string, file string) (*bufio.Scanner, error) {
	args := []string{dir, "-iname", file}
	cmd := exec.Command("find", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return bufio.NewScanner(&out), nil
}
