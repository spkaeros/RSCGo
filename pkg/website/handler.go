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
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/spkaeros/rscgo/pkg/server/log"
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

//Start Binds to the web port 8080 and serves HTTP content to it.
// Note: This is a blocking call, it will not return to caller.
func Start() {
	muxCtx.Handle("/", http.RedirectHandler("/404.html", 302))
	muxCtx.Handle("/index.ws", indexHandler())

	err := http.ListenAndServe(":8080", muxCtx)
	if err != nil {
		log.Error.Println("Could not bind to website port:", err)
		os.Exit(99)
	}
}
