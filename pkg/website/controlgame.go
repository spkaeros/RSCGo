/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package website

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/procexec"
	"github.com/spkaeros/rscgo/pkg/rand"
)

type buffers map[uint64]chan []byte

func bindGameProcManager() {
	var stdout io.Reader
	var bufs = make(buffers)
	var lockBufs sync.RWMutex
	var backBuf = make([][]byte, 0, 100)
	var ServerCmd *exec.Cmd
	var done = make(chan struct{})

	//muxCtx.Handle("/game/control.ws", pageHandler("Game Server Control", controlPage))
	muxCtx.HandleFunc("/game/launch.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd != nil && (ServerCmd.ProcessState == nil || !ServerCmd.ProcessState.Exited()) {
			writeContent(w, []byte("killing old instances...\n"))
			_ = procexec.Command("pkill", "-9", "game").Run()
		}

		ServerCmd = procexec.Run("game", "./bin/game", "-v")

		out, err := ServerCmd.StdoutPipe()
		if err != nil {
			writeContent(w, []byte("Error piping stdout:"+err.Error()))
			return
		}
		e, err := ServerCmd.StderrPipe()
		if err != nil {
			writeContent(w, []byte("Error piping stderr:"+err.Error()))
			return
		}
		stdout = io.MultiReader(out, e)

		err = ServerCmd.Start()
		if err != nil {
			if !writeContent(w, []byte("Error starting game server:"+err.Error())) {
				return
			}
		}
		go func() {
			err := ServerCmd.Wait()
			if err != nil && !strings.Contains(err.Error(), "killed") {
				log.Warning.Println("Error waiting for server command to finish running:", err)
				log.Warning.Printf("%v\n", ServerCmd.ProcessState)
				return
			}
			done <- struct{}{}
			if ServerCmd != nil && ServerCmd.ProcessState != nil {
				if failureCode := ServerCmd.ProcessState.ExitCode(); failureCode != 0 {
					log.Warn("Server exited with failure code:", failureCode)
					log.Warn(ServerCmd.ProcessState.String())
				}
			}
		}()
		go func() {
			b := bufio.NewReader(stdout)
			for ServerCmd != nil {
				line, err := b.ReadBytes('\n')
				if err != nil {
					return
				}
				backBuf = append(backBuf, line)
				if len(backBuf) == 100 {
					// truncates ten oldest buffered lines
					backBuf = backBuf[10:]
				}
				fmt.Printf("[GAME] %s", line)
				lockBufs.RLock()
				for _, buf := range bufs {
					buf <- line
				}
				lockBufs.RUnlock()
			}
		}()
		writeContent(w, []byte("Successfully started game server (pid: "+strconv.Itoa(ServerCmd.Process.Pid)+")"))
	})
	muxCtx.HandleFunc("/game/kill.ws", func(w http.ResponseWriter, r *http.Request) {
		pkill := func() error {
			cmd := procexec.Command("pkill", "game")
			err := cmd.Run()
			if err != nil {
				writeContent(w, []byte("Error killing the game server process:"+err.Error()))
			}
			return err
		}
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd == nil || ServerCmd.Process == nil || (ServerCmd.ProcessState != nil && ServerCmd.ProcessState.Exited()) {
			if err := pkill(); err != nil {
				writeContent(w, []byte("Error:['"+err.Error()+"]''; could not stop game server.  Is it running?"))
			}
			writeContent(w, []byte("Successfully killed game server"))
			return
		}
		err := ServerCmd.Process.Kill()
		ServerCmd = nil
		if err != nil {
			writeContent(w, []byte("Error:['"+err.Error()+"]''; Falling back to pkill..."))
			if err := pkill(); err != nil {
				writeContent(w, []byte("Error:['"+err.Error()+"]''; could not stop game server.  Is it running?"))
				return
			}
		}

		writeContent(w, []byte("Successfully killed game server"))
	})
	muxCtx.HandleFunc("/api/stdout", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Error.Printf("upgrade error: %s", err)
			return
		}
		defer conn.Close()
		for _, line := range backBuf {
			if err := wsutil.WriteServerText(conn, line); err != nil {
				log.Info.Println(err)
				return
			}
		}
		identifier := rand.Rng.Uint64()
		buf := make(chan []byte, 256)
		lockBufs.Lock()
		bufs[identifier] = buf
		lockBufs.Unlock()

		defer func() {
			lockBufs.Lock()
			delete(bufs, identifier)
			lockBufs.Unlock()
			close(buf)
		}()

		for {
			select {
			case line := <-buf:
				err := wsutil.WriteServerText(conn, line)
				if err != nil {
					return
				}
			}
		}
	})
}
