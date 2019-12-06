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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var Scripts []string

var EngineChannel = make(chan func(), 20)
//var InvTriggers []func(context.Context, reflect.Value, reflect.Value) (reflect.Value, reflect.Value)
//var BoundaryTriggers []func(context.Context, reflect.Value, reflect.Value) (reflect.Value, reflect.Value)

//var NpcTriggers []func(context.Context, reflect.Value, reflect.Value) (reflect.Value, reflect.Value)
var LoginTriggers []func(player *world.Player)
var InvOnBoundaryTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool
var InvOnObjectTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool
var InvTriggers = make(map[interface{}]func(player *world.Player, item *world.Item))
var ObjectTriggers = make(map[interface{}]func(*world.Player, *world.Object, int))
var BoundaryTriggers = make(map[interface{}]func(*world.Player, *world.Object, int))
var NpcTriggers = make(map[interface{}]func(*world.Player, *world.NPC))
var NpcAtkTriggers = make(map[interface{}]func(*world.Player, *world.NPC) bool)
var NpcDeathTriggers = make(map[interface{}]func(*world.Player, *world.NPC))

func Run(fnName string, player *world.Player, argName string, arg interface{}) bool {
	env := WorldModule()
	err := env.Define("client", player)
	if err != nil {
		log.Info.Println("Error initializing scripting environment:", err)
		return false
	}
	err = env.Define("player", player)
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
` + fnName + `()`)
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

func Clear() {
	//InvTriggers = InvTriggers[:0]
	InvTriggers = make(map[interface{}]func(*world.Player, *world.Item))
	//BoundaryTriggers = BoundaryTriggers[:0]
	ObjectTriggers = make(map[interface{}]func(*world.Player, *world.Object, int))
	BoundaryTriggers = make(map[interface{}]func(*world.Player, *world.Object, int))
	NpcTriggers = make(map[interface{}]func(*world.Player, *world.NPC))
	NpcDeathTriggers = make(map[interface{}]func(*world.Player, *world.NPC))
	NpcAtkTriggers = make(map[interface{}]func(*world.Player, *world.NPC) bool)
	LoginTriggers = LoginTriggers[:0]
	InvOnBoundaryTriggers = InvOnBoundaryTriggers[:0]
	InvOnObjectTriggers = InvOnObjectTriggers[:0]
}

//Load Loads all of the scripts in ./scripts and stores them in the Scripts slice.
func Load() {
	err := filepath.Walk("./scripts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Info.Println(err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, "ank") {
			env := WorldModule()
			_, err := env.Execute(load(path))
			if err != nil {
				log.Info.Println("Anko scripting error in '"+path+"':", err.Error())
				return nil
			}
		}
		return nil
	})
	if err != nil {
		log.Info.Println(err)
		return
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
