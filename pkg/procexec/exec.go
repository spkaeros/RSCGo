package procexec

import "os/exec"

func Command(command string, args ...string) *exec.Cmd {
	return exec.Command(command, args...)
}
