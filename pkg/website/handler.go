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
	"io"
	"fmt"
	"bufio"
	"bytes"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"strconv"

	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/gorilla/websocket"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var muxCtx = http.NewServeMux()

type InformationData struct {
	Title string
	Owner string
	Copyright string
}

var Information = InformationData{
	Title: "RSCGo",
	Owner: "ZlackCode LLC",
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

var StdOut, StdErr bytes.Buffer
var StdOutIn, StdErrIn io.Reader
var StdOutR, StdErrR io.Writer
var ServerProc *os.Process
//Start Binds to the web port 8080 and serves HTTP content to it.
// Note: This is a blocking call, it will not return to caller.
func Start() {
	muxCtx.Handle("/", http.NotFoundHandler())
	muxCtx.Handle("/index.ws", indexHandler())
	muxCtx.HandleFunc("/game/launch.ws", func(w http.ResponseWriter, r *http.Request) {
		if ServerProc != nil {
			w.Write([]byte("server already started\n"))
			return
		}
		cmd := exec.Command("./server", "-v")

		var err error
		StdOutIn, err = cmd.StdoutPipe()
		if err != nil {
			w.Write([]byte("Error setting stdout to var:" + err.Error()))
		}
		StdErrIn, err = cmd.StderrPipe()
		if err != nil {
			w.Write([]byte("Error setting stdout to var:" + err.Error()))
		}
		StdOutR = io.MultiWriter(os.Stdout, &StdOut)
		StdErrR = io.MultiWriter(os.Stderr, &StdErr)
		err = cmd.Start()
		if err != nil {
			w.Write([]byte("Error starting server process:" + err.Error()))
			return
		}
/*		go func() {
			_, err = io.Copy(StdOutR, StdOutIn,)
			if err != nil {
				w.Write([]byte("Error piping stdout to var:" + err.Error()))
			}
			_, err = io.Copy(StdErrR, StdErrIn,)
			if err != nil {
				w.Write([]byte("Error piping stderr to var:" + err.Error()))
			}
		}()
*/		ServerProc = cmd.Process
		w.Write(html)
		w.Write([]byte("Started server, process: " + strconv.Itoa(ServerProc.Pid) + "."))
	})
	muxCtx.HandleFunc("/game/shutdown.ws", func(w http.ResponseWriter, r *http.Request) {
		if ServerProc == nil {
			return
		}
		cmd := exec.Command("kill","-9",strconv.Itoa(ServerProc.Pid))
		err := cmd.Run()
		if err != nil {
			w.Write([]byte("Error starting server process:" + err.Error()))
			return
		}
		w.Write([]byte("gameserver killed successfully"))
		ServerProc = nil
	})
	muxCtx.HandleFunc("/game/out.ws", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(StdOut.Bytes()))
		w.Write([]byte(StdErr.Bytes()))
	})
//	ctx := http.NewServeMux()
	muxCtx.HandleFunc("/game/out", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte(StdOut.Bytes()))
//		w.Write([]byte(StdErr.Bytes()))
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

		s := bufio.NewScanner(io.MultiReader(StdOutIn,StdErrIn))
		for s.Scan() {
			ws.WriteMessage(1, s.Bytes())
		}
	})
	muxCtx.HandleFunc("/game/wot.ws", func(w http.ResponseWriter, r *http.Request) {
	})
/*	go func() {
		err := http.ListenAndServe(":8081", muxCtx)
		if err != nil {
			log.Error.Println("Could not bind to website port:", err)
			os.Exit(98)
		}
	}()
*/	err := http.ListenAndServe(":8080", muxCtx)
	if err != nil {
		log.Error.Println("Could not bind to website port:", err)
		os.Exit(99)
	}
}
