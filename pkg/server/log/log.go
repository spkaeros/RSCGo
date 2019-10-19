package log

import (
	"log"
	"os"
)

var (
	//Warning Log interface for warnings.
	Warning = log.New(os.Stdout, "[WARNING] ", log.Ltime|log.Lshortfile)
	//Info Log interface for debug information.
	Info = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	//Error Log interface for errors.
	Error = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	//Commands Log interface for in-game commands.
	Commands = log.New(os.Stdout, "[COMMAND] ", log.Ltime)
)
