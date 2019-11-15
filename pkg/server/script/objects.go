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
	"github.com/spkaeros/rscgo/pkg/server/log"
	"io/ioutil"
	"os"
)

var Scripts []string

//LoadObjectTriggers Loads all of the Tengo scripts in ./scripts/objects and stores them in the ObjectTriggers slice.
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

//LoadItemTriggers Loads all of the Tengo scripts in ./scripts/items and stores them in the ItemTriggers slice.
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
