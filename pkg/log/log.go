package log

import (
	"io"
	"log"
	"os"
)

var (
	//Info Log interface for debug information.
	//Info logs debug information.
	Info = log.New(os.Stdout, "[INFO] ", log.Ltime|log.Lshortfile)
	//Warning logs warning information.
	Warning = log.New(os.Stderr, "[WARNING] ", log.Ltime|log.Lshortfile)
	//Error logs error information.
	Error = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
	//Suspicious logs suspicious behavior.
	Suspicious = log.New(os.Stdout, "[SUSPICIOUS] ", log.Ltime)
	//Commands logs suspicious behavior.
	Commands = log.New(os.Stdout, "[COMMAND] ", log.Ltime)
)

func init() {
	dir := "." + string(os.PathSeparator) + "logs"
	if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		Error.Println("Error obtaining a directory to hold log files.  Using current working directory.", err)
		dir = dir[:1]
	}

	if outFile, err := os.OpenFile(dir+string(os.PathSeparator)+"out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		Info.SetOutput(io.MultiWriter(outFile, os.Stdout))
	} else {
		Error.Println("Could not open debug log file for writing:", err)
	}

	if outFile, err := os.OpenFile(dir+string(os.PathSeparator)+"warn.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		Warning.SetOutput(io.MultiWriter(outFile, os.Stderr))
	} else {
		Error.Println("Could not open warning log file for writing:", err)
	}

	if outFile, err := os.OpenFile(dir+string(os.PathSeparator)+"err.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		Warning.SetOutput(io.MultiWriter(outFile, os.Stderr))
	} else {
		Error.Println("Could not open error log file for writing:", err)
	}

	if outFile, err := os.OpenFile(dir+string(os.PathSeparator)+"cmd.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		Commands.SetOutput(io.MultiWriter(outFile, os.Stdout))
	} else {
		Error.Println("Could not open commands log file for writing:", err)
	}

	if outFile, err := os.OpenFile(dir+string(os.PathSeparator)+"cheaters.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		Suspicious.SetOutput(io.MultiWriter(outFile, os.Stdout))
	} else {
		Error.Println("Could not open cheaters log file for writing:", err)
	}
}

var Debugf = Info.Printf
var Debug = Info.Println
var Debugln = Info.Println

var Warnf = Warning.Printf
var Warn = Warning.Println

var Errorf = Error.Printf
var Fatal = Error.Println

var Cheatf = Suspicious.Printf
var Cheat = Suspicious.Println

var Command = Commands.Println
var Commandf = Commands.Printf
