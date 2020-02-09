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
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/procexec"
	"github.com/spkaeros/rscgo/pkg/rand"
)

var muxCtx = http.NewServeMux()

type InformationData struct {
	PageTitle     string
	Title     string
	Owner     string
	Copyright string
}

var Information = InformationData{
	PageTitle: "",
	Title:     "RSCGo",
	Owner:     "ZlackCode LLC",
	Copyright: "2019-2020",
}

func (s InformationData) ToLower(s2 string) string {
	return strings.ToLower(s2)
}

func (s InformationData) OnlineCount() int {
	return world.Players.Size()
}

var indexPage = template.Must(template.ParseFiles("./website/index.gohtml"))

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := indexPage.Execute(w, Information)
		if err != nil {
			log.Error.Println("Could not execute template template:", err)
			return
		}
	})
}

var stdout io.Reader
type buffers = map[uint64]chan []byte
type bufferSet struct{
	buffers
	sync.RWMutex
}
var outBuffers = bufferSet{buffers: make(buffers)}
var backBuffer = make([][]byte, 0, 1000)
var ServerCmd *exec.Cmd
var done = make(chan struct{})
var removing = make(chan int)

//writeContent is a helper function to write to a http.ResponseWriter easily with error handling
// returns true on success, otherwise false
func writeContent(w http.ResponseWriter, content []byte) bool {
	_, err := w.Write(content)
	if err != nil {
		log.Warning.Println("Error writing template to client:", err)
		return false
	}
	return true
}

func pageHandler(title string, template *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Information.PageTitle = title
		w.Header().Set("Content-Type", "text/html")
		err := template.ExecuteTemplate(w, "layout", Information)
		if err != nil {
			log.Error.Println("Could not execute layout template:", err)
			return
		}
	})
}

var controlPage = template.Must(template.ParseFiles("./website/layout.html", "./website/server_control.html"))

//Start Binds to the web port 8080 and serves HTTP template to it.
// Note: This is a blocking call, it will not return to caller.
func Start() {
	muxCtx.Handle("/", http.NotFoundHandler())
	muxCtx.Handle("/index.ws", indexHandler())
	muxCtx.Handle("/game/control.ws", pageHandler("Game Server Control", controlPage))
	muxCtx.HandleFunc("/game/launch.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd != nil {
			_, err := w.Write([]byte("game already started\n"))
			if err != nil {
				log.Warning.Println("Could not write game server control response:", err)
			}
			return
		}
		procexec.Command("pkill", "-9", "game").Run()
		
		ServerCmd = procexec.Run("rscgo", "./bin/game", "-v")

		out, err := ServerCmd.StdoutPipe()
		if err != nil {
			_, err := w.Write([]byte("Error making stdout pipe:" + err.Error()))
			if err != nil {
				log.Warning.Println("Could not write game server control response:", err)
			}
			return
		}
		e, err := ServerCmd.StderrPipe()
		if err != nil {
			_, err := w.Write([]byte("Error making stderr pipe:" + err.Error()))
			if err != nil {
				log.Warning.Println("Could not write game server control response:", err)
				return
			}
			return
		}
		stdout = io.MultiReader(out, e)

		err = ServerCmd.Start()
		if err != nil {
			_, err = w.Write([]byte("Error starting game server:"+err.Error()))
			if err != nil {
				log.Warning.Println("Could not write game server control response:", err)
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
				outBuffers.RLock()
				for _, buf := range outBuffers.buffers {
					buf <- line
				}
				outBuffers.RUnlock()
			}
		}()
		_, err = w.Write([]byte("Successfully started game server (pid: " + strconv.Itoa(ServerCmd.Process.Pid) + ")"))
		if err != nil {
			log.Warning.Println("Could not write game server control response:", err)
			return
		}
	})
	muxCtx.HandleFunc("/game/kill.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if ServerCmd == nil || ServerCmd.Process == nil || (ServerCmd.ProcessState != nil && ServerCmd.ProcessState.Exited()) {
			_, err := w.Write([]byte("Game server process could not be found.\n"))
			if err != nil {
				log.Warning.Println("Could not write game server control response:", err)
				return
			}
			return
		}
		err := ServerCmd.Process.Kill()
		if err != nil {
			cmd := procexec.Command("pkill", "game")
			err := cmd.Run()
			if err != nil {
				_, err := w.Write([]byte("Error killing the game server process:" + err.Error()))
				if err != nil {
					log.Warning.Println("Could not write game server control response:", err)
					return
				}
				return
			}
			return
		}
		_, err = w.Write([]byte("Successfully killed game server"))

		if err != nil {
			log.Warning.Println("Could not write game server control response:", err)
			return
		}
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
		outBuffers.Lock()
		outBuffers.buffers[identifier] = make(chan []byte, 256)
		outBuffers.Unlock()
		buf := outBuffers.buffers[identifier]

		defer func() {
			outBuffers.Lock()
			delete(outBuffers.buffers, identifier)
			outBuffers.Unlock()
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
	err := http.ListenAndServe(":8080", muxCtx)
	if err != nil {
		log.Error.Println("Could not bind to website port:", err)
		os.Exit(99)
	}
}
