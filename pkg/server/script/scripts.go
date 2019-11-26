/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package script

import (
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"io/ioutil"
	"os"
	"strings"
)

var Scripts []string

var TriggerC = make(chan func(), 20)

func Run(fnName string, c clients.Client, argName string, arg interface{}) bool {
	env := WorldModule()
	err := env.Define("client", c)
	if err != nil {
		log.Info.Println("Error initializing scripting environment:", err)
		return false
	}
	err = env.Define("player", c.Player())
	if err != nil {
		log.Info.Println("Error initializing scripting environment:", err)
		return false
	}
	err = env.Define(argName, arg)
	if err != nil {
		log.Info.Println("Error initializing scripting environment:", err)
		return false
	}
	for _, s := range Scripts {
		if !strings.Contains(s, fnName) {
			continue
		}
		stopPipeline, err := env.Execute(s +
			`
`+fnName+`()`)
		if err != nil {
			log.Info.Println("Unrecognized Anko error when attempting to execute the script pipeline:", err)
			continue
		}
		if stopPipeline, ok := stopPipeline.(bool); ok && stopPipeline {
			return true
		} else if !ok {
			log.Info.Println("Unexpected return result from an executed Anko script:", err)
		}
	}
	return false
}

//Load Loads all of the scripts in ./scripts and stores them in the Scripts slice.
func Load() {
	files, err := ioutil.ReadDir("./scripts")
	if err != nil {
		log.Info.Println("Error attempting to read scripts directory:", err)
		return
	}
	for _, file := range files {
		Scripts = append(Scripts, load("./scripts/" + file.Name()))
	}
}

func load(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Warning.Println("Error opening script file for object action:", err)
		return ""
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warning.Println("Error reading script file for object action:", err)
		return ""
	}

	return string(data)
}
