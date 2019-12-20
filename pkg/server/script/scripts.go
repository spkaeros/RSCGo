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
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mattn/anko/vm"

	"github.com/mattn/anko/parser"
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
)

var scriptWatcher *fsnotify.Watcher

//ItemTrigger A type that defines a callback to run when certain item actions are performed, and a predicate to decide
// whether or not the callback should run
type ItemTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.Item) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.Item)
}

//ObjectTrigger A type that defines a callback to run when certain object actions are performed, and a predicate to decide
// whether or not the callback should run
type ObjectTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.Object, int) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.Object, int)
}

//NpcTrigger A type that defines a callback to run when certain NPC actions are performed, and a predicate to decide
// whether or not the callback should run
type NpcTrigger struct {
	// Check returns true if this handler should run.
	Check func(*world.NPC) bool
	// Action is the function that will run if Check returned true.
	Action func(*world.Player, *world.NPC)
}

//NpcActionPredicate A type alias for an NPC related action predicate.
type NpcActionPredicate = func(*world.Player, *world.NPC) bool

//NpcAction A type alias for an NPC related action.
type NpcAction = func(*world.Player, *world.NPC)

//LoginTriggers a list of actions to run when a player logs in.
var LoginTriggers []func(player *world.Player)

//InvOnBoundaryTriggers a list of actions to run when a player uses an inventory item on a boundary object
var InvOnBoundaryTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool

//InvOnObjectTriggers a list of actions to run when a player uses an inventory item on a object
var InvOnObjectTriggers []func(player *world.Player, object *world.Object, item *world.Item) bool

//ItemTriggers List of script callbacks to run for inventory item actions
var ItemTriggers []ItemTrigger

//ObjectTriggers List of script callbacks to run for object actions
var ObjectTriggers []ObjectTrigger

//BoundaryTriggers List of script callbacks to run for boundary actions
var BoundaryTriggers []ObjectTrigger

//NpcTriggers List of script callbacks to run for NPC talking actions
var NpcTriggers []NpcTrigger

//NpcAtkTriggers List of script callbacks to run when you attack an NPC
var NpcAtkTriggers []world.NpcBlockingTrigger

//Clear clears all of the lists of triggers.
func Clear() {
	ItemTriggers = ItemTriggers[:0]
	ObjectTriggers = ObjectTriggers[:0]
	NpcTriggers = NpcTriggers[:0]
	NpcAtkTriggers = NpcAtkTriggers[:0]
	world.NpcDeathTriggers = world.NpcDeathTriggers[:0]
	BoundaryTriggers = BoundaryTriggers[:0]
	LoginTriggers = LoginTriggers[:0]
	InvOnBoundaryTriggers = InvOnBoundaryTriggers[:0]
	InvOnObjectTriggers = InvOnObjectTriggers[:0]
}

//Load Loads all of the scripts in ./scripts.  This will ignore any folders named definitions or lib.
func Load() {
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
						_, err := vm.Execute(WorldModule(), &vm.Options{Debug: true}, load(event.Name))

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
	parser.EnableDebug(1)
	parser.EnableErrorVerbose()

	err = filepath.Walk("./scripts", func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsDir() && !strings.Contains(path, "definitions") && !strings.Contains(path, "lib") && strings.HasSuffix(path, "ank") {
			_, err := vm.Execute(WorldModule(), &vm.Options{Debug: true}, load(path))

			if err != nil {
				log.Info.Println("Anko error ['"+path+"']:", err)
				//				log.Info.Println(env.String())
				return nil
			}
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
		log.Warning.Println("Error opening script file:", err)
		return ""
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warning.Println("Error reading script file:", err)
		return ""
	}

	return string(data)
}
