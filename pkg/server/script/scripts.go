package script

import (
	"bitbucket.org/zlacki/rscgo/pkg/server/log"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)


//Initialize Initializes a Tengo script with the specified data.
func Initialize(data string) *script.Script {
	s := script.New([]byte(data))
	scriptModules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	scriptModules.Remove("os")
	scriptModules.Add("world", NewWorldModule())
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
