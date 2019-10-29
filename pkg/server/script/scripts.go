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

	return Initialize(string(data))
}

//Initialize Initializes a Tengo script with the specified data.
func Initialize(data string) *script.Script {
	s := script.New([]byte(data))
 	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", world.NewWorldModule())
	SetScriptVariable(s, "ret", objects.FalseValue)
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