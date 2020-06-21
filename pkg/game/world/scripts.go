/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package world

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mattn/anko/vm"

	"github.com/spkaeros/rscgo/pkg/log"
)

var scriptWatcher *fsnotify.Watcher

//ItemTrigger A type that defines a callback to run when certain item actions are performed, and a predicate to decide
// whether or not the callback should run
type ItemTrigger struct {
	// Check returns true if this handler should run.
	Check func(*Item) bool
	// Action is the function that will run if Check returned true.
	Action func(*Player, *Item)
}

//ItemOnPlayerTrigger A type that defines a callback to run when certain item are used on players, and a predicate to decide
// whether or not the callback should run
type ItemOnPlayerTrigger struct {
	// Check returns true if this handler should run.
	Check func(*Item) bool
	// Action is the function that will run if Check returned true.
	Action func(*Player, *Player, *Item)
}

//ObjectTrigger A type that defines a callback to run when certain object actions are performed, and a predicate to decide
// whether or not the callback should run
type ObjectTrigger struct {
	// Check returns true if this handler should run.
	Check func(*Object, int) bool
	// Action is the function that will run if Check returned true.
	Action func(*Player, *Object, int)
}

//NpcTrigger A type that defines a callback to run when certain NPC actions are performed, and a predicate to decide
// whether or not the callback should run
type NpcTrigger struct {
	// Check returns true if this handler should run.
	Check func(*NPC) bool
	// Action is the function that will run if Check returned true.
	Action func(*Player, *NPC)
}

type Trigger func(*Player, interface{})

//NpcActionPredicate A type alias for an NPC related action predicate.
type NpcActionPredicate = func(*Player, *NPC) bool

//NpcAction A type alias for an NPC related action.
type NpcAction = func(*Player, *NPC)

//LoginTriggers a list of actions to run when a player logs in.
var LoginTriggers []func(player *Player)

//InvOnBoundaryTriggers a list of actions to run when a player uses an inventory item on a boundary object
var InvOnBoundaryTriggers []func(player *Player, object *Object, item *Item) bool

//InvOnObjectTriggers a list of actions to run when a player uses an inventory item on a object
var InvOnObjectTriggers []func(player *Player, object *Object, item *Item) bool

//ItemTriggers List of script callbacks to run for inventory item actions
var ItemTriggers []ItemTrigger

//InvOnPlayerTriggers a list of actions to run when a player uses an inventory item on another player object
var InvOnPlayerTriggers []ItemOnPlayerTrigger

//ObjectTriggers List of script callbacks to run for object actions
var ObjectTriggers []ObjectTrigger

//BoundaryTriggers List of script callbacks to run for boundary actions
var BoundaryTriggers []ObjectTrigger

//NpcTriggers List of script callbacks to run for NPC talking actions
var NpcTriggers []NpcTrigger

//var Triggers []Trigger

var SpellTriggers = make(map[int]Trigger)

type SpellDef map[string]interface{}

//NpcAtkTriggers List of script callbacks to run when you attack an NPC
var NpcAtkTriggers []NpcBlockingTrigger

//Clear clears all of the lists of triggers.
func Clear() {
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

//RunScripts Loads all of the scripts in ./scripts.  This will ignore any folders named definitions or lib.
func RunScripts() {
	var err error
	if scriptWatcher == nil {
		scriptWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Info.Println(err)
			return
		}
		go func() {
			lastEvent := time.Now()
			lastPath := ""
			for {
				select {
				case event := <-scriptWatcher.Events:
					if time.Since(lastEvent) < time.Second && lastPath == event.Name {
						continue
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						lastEvent = time.Now()
						lastPath = event.Name
						log.Info.Println("Reloading " + event.Name)
						_, err := vm.Execute(ScriptEnv(), &vm.Options{Debug: true}, load(event.Name))

						if err != nil {
							log.Info.Println("Anko error ['"+event.Name+"']:", err)
							continue
						}
					}
				case err := <-scriptWatcher.Errors:
					if err != nil {
						log.Info.Println(err)
					}
				}
			}
		}()
	}

	err = filepath.Walk("./scripts", func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsDir() && strings.HasSuffix(path, "ank") && !strings.Contains(path, "defyyyyyyyt6") && !strings.Contains(path, "lib") {

			_, err := vm.Execute(ScriptEnv(), &vm.Options{Debug: true}, "bind = import(\"bind\")\nworld = import(\"world\")\nlog = import(\"log\")\nids = import(\"ids\")\n\n"+load(path))
			//stmt, err := parser.ParseSrc(load(path))
			//if err != nil {
			//	log.Warning.Printf("ParseSrc error - received: %v - script: %v", err, path)
			//}
			//// Note: Still want to run the code even after a parse error to see what happens
			//_, err = vm.Run(ScriptEnv(), &vm.Options{Debug: true}, stmt)
			if err != nil {
				log.Warn("Anko error ['"+path+"']:", err)
				//				log.Info.Println(env.String())
				return nil
			}
			//log.Info.Println(val)
			return scriptWatcher.Add(path)
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
		log.Warn("Error opening script file:", err)
		return ""
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warn("Error reading script file:", err)
		return ""
	}

	return string(data)
}
