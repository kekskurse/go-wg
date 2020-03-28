package wg

import (
	"log"
	"os/exec"
)

type WGShellInterface interface {
	Command(name string, arg ...string) *exec.Cmd
}

type WGShell struct {

}

func (w WGShell) Command(name string, arg ...string) *exec.Cmd {
	completeCommand := name
	for _, v := range arg {
		completeCommand = completeCommand + " " + v
	}
	log.Println("Shell-Execute: ", completeCommand)
	return exec.Command(name, arg...)
}