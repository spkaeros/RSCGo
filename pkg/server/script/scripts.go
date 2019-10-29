package script

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/clients"
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"bitbucket.org/zlacki/rscgo/pkg/server/world"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"io/ioutil"
	"os"
)

var objectWrapper = []byte(`ret := import("main")(player, object, cmd)`)

var scriptAttributes = map[string]objects.Object {
	"replaceObjectAt": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 3 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			object := world.GetObject(x, y)
			if object == nil {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "GetObject(x,y)",
					Expected: "An object",
					Found:    "Nothing",
				}
			}
			id, ok := objects.ToInt(args[2])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "objectID",
					Expected: "int",
					Found:    args[2].TypeName(),
				}
			}
			world.ReplaceObject(object, id)
			return objects.UndefinedValue, nil
		},
	},
	"replaceObject": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			object, ok := args[0].(*world.Object)
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "object",
					Expected: "*world.Object",
					Found:    args[0].TypeName(),
				}
			}
			id, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "objectID",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			world.ReplaceObject(object, id)
			return objects.UndefinedValue, nil
		},
	},
	"getObject": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			object := world.GetObject(x, y)
			if object == nil {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "GetObject(x,y)",
					Expected: "An object",
					Found:    "Nothing",
				}
			}
			return object, nil
		},
	},
	"addObject": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 3 {
				return nil, objects.ErrWrongNumArguments
			}
			id, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "id",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			x, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[2])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[2].TypeName(),
				}
			}
			if object := world.GetObject(x, y); object != nil {
				world.ReplaceObject(object, id)
			} else {
				world.AddObject(world.NewObject(id, 0, x, y, false))
			}
			return objects.UndefinedValue, nil
		},
	},
	"removeObject": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			if object := world.GetObject(x, y); object == nil {
				world.RemoveObject(object)
			}
			return objects.UndefinedValue, nil
		},
	},
	"getPlayer": &objects.UserFunction {
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 1 {
				return nil, objects.ErrWrongNumArguments
			}
			index, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "index",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			client, ok := clients.FromIndex(index)
			if !ok {
				return nil, objects.ErrIndexOutOfBounds
			}
			return client, nil
		},
	},
}

func NewWorldModule() *objects.BuiltinModule {
	return &objects.BuiltinModule{ Attrs: scriptAttributes}
}

var ObjectTriggers []*script.Script

func LoadObjectTriggers() {
	files, err := ioutil.ReadDir("./scripts/objects")
	if err != nil {
		log.Info.Println("Error attempting to read scripts directory:", err)
		return
	}
	for _, file := range files {
		ObjectTriggers = append(ObjectTriggers, LoadObjectTrigger("./scripts/objects/" + file.Name()))
	}
}

//LoadObjectTrigger Loads the data in the file located at filePath on the local file system, and initializes a new Tengo VM script with it.
func LoadObjectTrigger(filePath string) *script.Script {
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

	return InitializeObjectTrigger(data)
}

//Initialize Initializes a Tengo script with the specified data.
func Initialize(data string) *script.Script {
	s := script.New([]byte(data))
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", NewWorldModule())
	s.SetImports(scriptModules)
	return s
}

//Initialize Initializes a Tengo script with the specified data, using a wrapper for object action triggers.
func InitializeObjectTrigger(data []byte) *script.Script {
	s := script.New(objectWrapper)
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", NewWorldModule())
	scriptModules.AddSourceModule("main", []byte(data))
	s.SetImports(scriptModules)
	return s
}

//RunScript Runs a script on the Tengo VM, with error checking.
func RunScript(s *script.Script) bool {
	compiled, err := s.Compile()
	if err != nil {
		log.Info.Println(err)
		return false
	}
	if err := compiled.Run(); err != nil {
		log.Info.Println("Error running script in VM:", err)
		return false
	}
	return compiled.Get("ret").Bool()
}

//SetScriptVariable Sets a script-scoped variable by name to value.
func SetScriptVariable(s *script.Script, variableName string, value interface{}) {
	if err := s.Add(variableName, value); err != nil {
		log.Info.Println("Error setting script variable '" + variableName + "':", err)
		return
	}
}

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