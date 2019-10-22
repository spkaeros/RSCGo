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
	//Suspicious Log interface for in-game suspicious behavior.
	Suspicious = log.New(os.Stdout, "[SUSPICIOUS] ", log.Ltime)
)

func init() {
	if f, err := os.OpenFile("logs"+string(os.PathSeparator)+"cmd.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		Error.Println("Could not open commands log file for writing:", err)
	} else {
		Commands.SetOutput(f)
	}
	if f, err := os.OpenFile("logs"+string(os.PathSeparator)+"cheaters.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		Error.Println("Could not open cheaters log file for writing:", err)
	} else {
		Suspicious.SetOutput(f)
	}
}