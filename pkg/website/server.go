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
	"path/filepath"
	"strings"

	"github.com/spkaeros/rscgo/pkg/game/world"
	"github.com/spkaeros/rscgo/pkg/log"
)

var muxCtx = http.NewServeMux()

type InformationData struct {
	PageTitle string
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

var templates = make(map[string]*template.Template)

// Load templates on program initialisation
func init() {
	layouts, err := filepath.Glob("website/*/*.html")
	if err != nil {
		log.Error.Fatal(err)
	}
	layouts2, err := filepath.Glob("website/*.html")

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range append(layouts, layouts2...) {
		templates[layout[8:]] = template.Must(template.ParseFiles("website/layouts/layout.html", layout))
	}
}

func render(w http.ResponseWriter, r *http.Request) {
	name := strings.ReplaceAll(filepath.Clean(r.URL.Path[1:]), ".ws", ".html")
	tmpl, ok := templates[name]
	if !ok {
		w.WriteHeader(404)
		writeContent(w, []byte("404 file not found"))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, "layout", Information)
	if err != nil {
		w.WriteHeader(500)
		writeContent(w, []byte("Internal Server Error"))
	}
}

//var controlPage = template.Must(template.ParseFiles("./website/layouts/layout.html", "./website/control.html"))

//Start Binds to the web port 8080 and serves HTTP template to it.
// Note: This is a blocking call, it will not return to caller.
func Start() {
	muxCtx.HandleFunc("/", render)
	muxCtx.HandleFunc("/game/", render)
	bindGameProcManager()
	err := http.ListenAndServe(":8080", muxCtx)
	if err != nil {
		log.Error.Println("Could not bind to website port:", err)
		os.Exit(99)
	}
}
