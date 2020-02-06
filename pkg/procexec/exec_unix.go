//+build !windows

/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package procexec

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spkaeros/rscgo/pkg/log"
)

//Command Attempts to make a *exec.Cmd setup to run the file at path, and tries system-specific
// API calls to prevent tying the child processes lifetime to the parents.
// Returns a valid *exec.Cmd that can run the file at path upon success, or nil upon failure
func Command(path string, name string) *exec.Cmd {
	s, err := os.Stat(path)
	if err != nil {
		log.Warning.Println("Could not stat file:", err)
		return nil
	}

	if s.IsDir() {
		log.Warning.Println("File at path '" + path + "' is not an executable binary file!")
		return nil

	}
	return getCmd(path)
}

func getCmd(path string) *exec.Cmd {
	cmd := exec.Command(path, "-v")
	cmd.Args[0] = "rscgo"
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	return cmd
}
