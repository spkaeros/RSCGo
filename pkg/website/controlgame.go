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

type buffers = map[uint64]chan []byte
type bufferSet struct{
	buffers
	sync.RWMutex
}

func bindGameProcManager() {
	var stdout io.Reader
	var stdoutClients = bufferSet{buffers: make(buffers)}
	var backBuffer = make([][]byte, 0, 1000)
	var ServerCmd *exec.Cmd
	var done = make(chan struct{})

	//muxCtx.Handle("/game/control.ws", pageHandler("Game Server Control", controlPage))
	muxCtx.HandleFunc("/game/launch.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd != nil {
			writeContent(w, []byte("game already started\n"))
			return
		}
		_ = procexec.Command("pkill", "-9", "game").Run()

		ServerCmd = procexec.Run("rscgo", "./bin/game", "-v")

		out, err := ServerCmd.StdoutPipe()
		if err != nil {
			writeContent(w, []byte("Error making stdout pipe:" + err.Error()))
			return
		}
		e, err := ServerCmd.StderrPipe()
		if err != nil {
			writeContent(w, []byte("Error making stderr pipe:" + err.Error()))
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
			done <- struct {}{}
			if ServerCmd != nil && ServerCmd.ProcessState != nil {
				if failureCode := ServerCmd.ProcessState.ExitCode(); failureCode != 0 {
					log.Warning.Println("Server exited with failure code:", failureCode)
					log.Warning.Println(ServerCmd.ProcessState.String())
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
				backBuffer = append(backBuffer, line)
				if len(backBuffer) == 1000 {
					backBuffer = backBuffer[100:]
				}
				fmt.Printf("[GAME] %s", line)
				stdoutClients.RLock()
				for _, buf := range stdoutClients.buffers {
					buf <- line
				}
				stdoutClients.RUnlock()
			}
		}()
		writeContent(w, []byte("Successfully started game server (pid: " + strconv.Itoa(ServerCmd.Process.Pid) + ")"))
	})
	muxCtx.HandleFunc("/game/kill.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd == nil || ServerCmd.Process == nil || (ServerCmd.ProcessState != nil && ServerCmd.ProcessState.Exited()) {
			writeContent(w, []byte("Game server process could not be found.\n"))
			return
		}
		err := ServerCmd.Process.Kill()
		if err != nil {
			cmd := procexec.Command("pkill", "game")
			err := cmd.Run()
			if err != nil {
				writeContent(w, []byte("Error killing the game server process:" + err.Error()))
				return
			}
			return
		}

		writeContent(w, []byte("Successfully killed game server"))
		ServerCmd = nil
	})
	muxCtx.HandleFunc("/api/game/stdout", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Error.Printf("upgrade error: %s", err)
			return
		}
		defer conn.Close()
		for _, line := range backBuffer {
			if err := wsutil.WriteServerText(conn, line); err != nil {
				log.Info.Println(err)
				return
			}
		}
		identifier := rand.Uint64()
		stdoutClients.Lock()
		stdoutClients.buffers[identifier] = make(chan []byte, 256)
		stdoutClients.Unlock()
		buf := stdoutClients.buffers[identifier]

		defer func() {
			stdoutClients.Lock()
			delete(stdoutClients.buffers, identifier)
			stdoutClients.Unlock()
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
