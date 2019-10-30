package script

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
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
var ObjectTriggers []*script.Script

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

	return InitializeObjectTrigger(append(data, objectWrapper...))
}

//InitializeObjectTrigger Initializes a Tengo script with the specified data, using a wrapper for object action triggers.
func InitializeObjectTrigger(data []byte) *script.Script {
	s := script.New(data)
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", NewWorldModule())
//	scriptModules.AddSourceModule("main", []byte(data))
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