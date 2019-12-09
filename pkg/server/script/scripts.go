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

	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var Scripts []string

//ItemTrigger holds callbacks to functions defined in the Anko scripts loaded at runtime, to be run when certain
// events occur
type ItemTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.Item) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.Item)
}

type ObjectTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.Object, int) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.Object, int)
}

//NpcTrigger holds callbacks to functions defined in the Anko scripts loaded at runtime, to be run when certain
// events occur
type NpcTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.NPC) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.NPC)
}

//NpcActionPredicate callback to a function defined in the Anko scripts loaded at runtime, to be run when certain
// events occur.  If it returns true, it will block the event that triggered it from occurring
type NpcPredBlockingTrigger struct {
	// Check returns true if this handler should run.
	Check NpcActionPredicate
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.NPC)
}

type NpcActionPredicate = func(*world.Player, *world.NPC) bool
type NpcAction = func(*world.Player, *world.NPC)

var LoginTriggers []func(player *world.Player)
var InvOnBoundaryTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool
var InvOnObjectTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool

//ItemTriggers List of script callbacks to run for inventory item actions
var ItemTriggers []ItemTrigger
var ObjectTriggers []ObjectTrigger
var BoundaryTriggers []ObjectTrigger

//NpcTriggers List of script callbacks to run for NPC talking actions
var NpcTriggers []NpcTrigger

//NpcAtkTriggers List of script callbacks to run when you attack an NPC
var NpcAtkTriggers []NpcActionPredicate

//NpcDeathTriggers List of script callbacks to run when you kill an NPC
var NpcDeathTriggers []NpcPredBlockingTrigger

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
		stopPipeline, err := vm.Execute(env, nil, s+
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

func Clear() {
	//ItemTriggers = make(map[interface{}]func(*world.Player, *world.Item))
	ItemTriggers = ItemTriggers[:0]
	ObjectTriggers = ObjectTriggers[:0]
	NpcTriggers = NpcTriggers[:0]
	NpcAtkTriggers = NpcAtkTriggers[:0]
	NpcDeathTriggers = NpcDeathTriggers[:0]
	BoundaryTriggers = BoundaryTriggers[:0]
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
		if !info.IsDir() && !strings.Contains(path, "definitions") && !strings.Contains(path, "lib") && strings.HasSuffix(path, "ank") {
			env := WorldModule()
			parser.EnableDebug(1)
			parser.EnableErrorVerbose()
			_, err := vm.Execute(env, &vm.Options{Debug: true}, load(path))

			if err != nil {
				log.Info.Println("Anko scripting error in '"+path+"':", err)
				//				log.Info.Println(env.String())
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
