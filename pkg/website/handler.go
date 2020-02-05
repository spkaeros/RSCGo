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
	"time"

	"github.com/gorilla/websocket"
	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

var muxCtx = http.NewServeMux()

type InformationData struct {
	Title     string
	Owner     string
	Copyright string
}

var Information = InformationData{
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

var index = template.Must(template.ParseFiles("./website/index.gohtml"))

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := index.Execute(w, Information)
		if err != nil {
			log.Error.Println("Could not execute html template:", err)
			return
		}
	})
}

var html = []byte(
	`<html>
	<body>
		<h1>game process stdout/stderr</h1>
		<code></code>
		<script>
			var ws = new WebSocket("wss://rscturmoil.com/game/out")
			ws.onmessage = function(e) {
				document.querySelector("code").innerHTML += e.data + "<br>"
			}
		</script>
	</body>
</html>
`)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var outBuffer = make(chan []byte, 256)
var ServerProc *os.Process

//Start Binds to the web port 8080 and serves HTTP content to it.
// Note: This is a blocking call, it will not return to caller.
func Start() {
	muxCtx.Handle("/", http.NotFoundHandler())
	muxCtx.Handle("/index.ws", indexHandler())
	muxCtx.HandleFunc("/game/launch.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
		if ServerProc != nil {
			w.Write([]byte("game already started\n"))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		cmd := exec.Command("./game", "-v")

		outReader, err := cmd.StdoutPipe()
		if err != nil {
			w.Write([]byte("Error getting game game output pipe reader:" + err.Error()))
		}
		errReader, err := cmd.StderrPipe()
		if err != nil {
			w.Write([]byte("Error getting game game error pipe reader:" + err.Error()))
		}
		scanner := bufio.NewScanner(io.MultiReader(outReader, errReader))
		//		multiWriter := io.MultiWriter(os.Stdout, &outBuffer)
		go func() {
			for scanner.Scan() {
				os.Stdout.Write(append(scanner.Bytes(), byte('\n')))
				outBuffer <- []byte(scanner.Text())
			}
		}()
		err = cmd.Start()
		if err != nil {
			w.Write([]byte("Error starting game process:" + err.Error()))
			return
		}
		ServerProc = cmd.Process
		w.Write([]byte("Started game, process: " + strconv.Itoa(ServerProc.Pid) + "."))
	})
	muxCtx.HandleFunc("/game/shutdown.ws", func(w http.ResponseWriter, r *http.Request) {
		if ServerProc == nil {
			w.Write([]byte("game child process not launched.\n"))
			return
		}
		cmd := exec.Command("kill", "-9", strconv.Itoa(ServerProc.Pid))
		err := cmd.Run()
		if err != nil {
			w.Write([]byte("Error starting kill process:" + err.Error()))
			return
		}
		w.Write([]byte("Game game(" + strconv.Itoa(ServerProc.Pid) + ") shut down successfully"))
		ServerProc = nil
	})
	muxCtx.HandleFunc("/game/out.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Write(html)
	})
	muxCtx.HandleFunc("/game/out", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("", err)))
			return
		}
		defer ws.Close()
		go func(c *websocket.Conn) {
			for {
				if _, _, err := c.NextReader(); err != nil {
					c.Close()
					break
				}
			}
		}(ws)

		ticker := time.NewTicker(54 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case line, ok := <-outBuffer:
				ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if !ok {
					ws.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				if err := ws.WriteMessage(1, line); err != nil {
					log.Error.Println(err)
					return
				}
			case <-ticker.C:
				ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Info.Println(err)
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
