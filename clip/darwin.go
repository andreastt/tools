// +build darwin

package main

import (
	"io"
	"os/exec"
	"bytes"
	"bufio"
)

func copy(r io.Reader) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = r
	return cmd.Run()
}

func paste() (io.Reader, error) {
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return bufio.NewReader(&out), nil
}
