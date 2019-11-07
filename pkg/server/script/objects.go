package script

import (
	"github.com/spkaeros/rscgo/pkg/server/log"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"io/ioutil"
	"os"
)

var objectWrapper = []byte(`
ret := isHandled(player, object, cmd)
if ret {
	action(player, object, cmd)
}`)
var itemWrapper = []byte(`
ret := isHandled(player, item, cmd)
if ret {
	action(player, item, cmd)
}`)
var ObjectTriggers []*script.Script
var ItemTriggers []*script.Script
var BoundaryTriggers []*script.Script

//LoadObjectTriggers Loads all of the Tengo scripts in ./scripts/objects and stores them in the ObjectTriggers slice.
func LoadObjectTriggers() {
	files, err := ioutil.ReadDir("./scripts/objects")
	if err != nil {
		log.Info.Println("Error attempting to read scripts directory:", err)
		return
	}
	for _, file := range files {
		ObjectTriggers = append(ObjectTriggers, loadObjectTrigger("./scripts/objects/" + file.Name()))
	}
}

//LoadItemTriggers Loads all of the Tengo scripts in ./scripts/items and stores them in the ItemTriggers slice.
func LoadItemTriggers() {
	files, err := ioutil.ReadDir("./scripts/items")
	if err != nil {
		log.Info.Println("Error attempting to read scripts directory:", err)
		return
	}
	for _, file := range files {
		ItemTriggers = append(ItemTriggers, loadItemTrigger("./scripts/items/" + file.Name()))
	}
}

//LoadObjectTriggers Loads all of the Tengo scripts in ./scripts/objects and stores them in the ObjectTriggers slice.
func LoadBoundaryTriggers() {
	files, err := ioutil.ReadDir("./scripts/boundarys")
	if err != nil {
		log.Info.Println("Error attempting to read scripts directory:", err)
		return
	}
	for _, file := range files {
		BoundaryTriggers = append(BoundaryTriggers, loadBoundaryTrigger("./scripts/boundarys/" + file.Name()))
	}
}

//loadObjectTrigger Loads the data in the file located at filePath on the local file system, and initializes a new Tengo VM script with it.
func loadObjectTrigger(filePath string) *script.Script {
	file, err := os.Open(filePath)
	if err != nil {
		log.Warning.Println("Error opening script file for object action:", err)
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warning.Println("Error reading script file for object action:", err)
		return nil
	}

	return initializeTrigger(append(data, objectWrapper...))
}

//loadItemTrigger Loads the data in the file located at filePath on the local file system, and initializes a new Tengo VM script with it.
func loadItemTrigger(filePath string) *script.Script {
	file, err := os.Open(filePath)
	if err != nil {
		log.Warning.Println("Error opening script file for item action:", err)
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warning.Println("Error reading script file for item action:", err)
		return nil
	}

	return initializeTrigger(append(data, itemWrapper...))
}

//loadBoundaryTrigger Loads the data in the file located at filePath on the local file system, and initializes a new Tengo VM script with it.
func loadBoundaryTrigger(filePath string) *script.Script {
	file, err := os.Open(filePath)
	if err != nil {
		log.Warning.Println("Error opening script file for boundary action:", err)
		return nil
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warning.Println("Error reading script file for boundary action:", err)
		return nil
	}

	return initializeTrigger(append(data, objectWrapper...))
}

//initializeTrigger Initializes a Tengo script with the specified data, using a wrapper for object action triggers.
func initializeTrigger(data []byte) *script.Script {
	s := script.New(data)
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", NewWorldModule())
	s.SetImports(scriptModules)
	return s
}

//ReplaceObject Why is this here?
func ReplaceObject(o *world.Object) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			ret = objects.UndefinedValue
			if len(args) < 1 {
				return nil, objects.ErrWrongNumArguments
			}
			id, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "objectID",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			world.ReplaceObject(o, id)
			return
		},
	}
}