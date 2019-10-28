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

var worldScript = []byte(`
fmt := import("fmt")

export (
	
)
`)

//Load Loads the data in the file located at filePath on the local file system, and initializes a new Tengo VM script with it.
func Load(filePath string) *script.Script {
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

	s := script.New(data)
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	s.SetImports(scriptModules)
//	SetScriptVariable(s, "ret", "")
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

func ScriptMessage(c clients.Client) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			ret = objects.UndefinedValue

			message, ok := objects.ToString(args[0])
			if !ok {
				message = args[0].String()
			}

			c.Message(message)
			return
		},
	}
}

func MovePlayer(c clients.Client) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			ret = objects.UndefinedValue
			if len(args) != 2 {
				c.Message("teleport(x,y): Invalid argument count provided")
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				c.Message("teleport(x,y): Invalid argument type provided")
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				c.Message("teleport(x,y): Invalid argument type provided")
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			c.Player().Teleport(x, y)
			return
		},
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

func ClimbUp(c clients.Client) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			ret = objects.UndefinedValue
			if nextLocation := c.Player().Above(); !nextLocation.Equals(c.Player().Location) {
				c.Player().ResetPath()
				c.Player().SetLocation(&nextLocation)
				c.UpdatePlane()
			}
			return
		},
	}
}

func ClimbDown(c clients.Client) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			ret = objects.UndefinedValue
			if nextLocation := c.Player().Below(); !nextLocation.Equals(c.Player().Location) {
				c.Player().ResetPath()
				c.Player().SetLocation(&nextLocation)
				c.UpdatePlane()
			}
			return
		},
	}
}