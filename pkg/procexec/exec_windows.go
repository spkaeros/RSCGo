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

//Run Attempts to make a *exec.Cmd setup to run the specified file, with system-specific
// API calls to prevent tying the child processes lifetime to the parents.
// Returns a valid *exec.Cmd that can run the file upon success, or nil upon failure
func Run(name string, file string, args ...string) *exec.Cmd {
	s, err := os.Stat(file)
	if err != nil {
		log.Warn("Could not stat file:", err)
		return nil
	}

	if !s.Mode().IsRegular() {
		log.Warn("File at path '" + file + "' is not an executable file!")
		return nil
	}

	cmd := exec.Command(file, args...)
	cmd.Args[0] = name + ".exe"
	// HideWindow to prevent a sudden console popping up on the host machine
	// CreationFlags to prevent cmd.Process dying when the caller process dies
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

	return cmd
}
