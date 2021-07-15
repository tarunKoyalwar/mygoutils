package myio

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
)

// A commadData is struct to return data
// after command is run sucessfully
// Note stderr is not returned
type CommandData struct {
	Stdout  []byte      // A byte array containing program output
	Process *os.Process // Process pointer to get details
	Name    string      // Command name  or path
}

// A function used to check if command exists and check path
// In future version will implement command installations
func CheckPath(command string) string {
	cmddata, err := exec.LookPath(command)
	HandleError(err, "Command Not found please install")
	return cmddata
}

// This function is used to run commands with additional options
// Additional options such as stderr , current directory etc
// It takes CMD struct as input
func RuncmdAdv(ptr *exec.Cmd) *CommandData {
	var out bytes.Buffer

	ptr.Stdout = &out

	err := ptr.Run()

	HandleError(err, "Command Failed to start => ")

	return &CommandData{
		Stdout:  out.Bytes(),
		Name:    ptr.Path,
		Process: ptr.Process,
	}

}

// This function is used to run commands with additional options
// Additional options such as stderr , current directory etc
// It takes CMD struct as input
// useful for goroutines
func GorunAdvcmd(ptr *exec.Cmd, ch chan<- CommandData, wg *sync.WaitGroup) {
	var out bytes.Buffer

	ptr.Stdout = &out

	err := ptr.Run()

	HandleError(err, "Command Failed to start => ")

	datax := CommandData{
		Stdout:  out.Bytes(),
		Name:    ptr.Path,
		Process: ptr.Process,
	}

	ch <- datax

	defer wg.Done()

}

// A simple function to run command and return its data
// takes cmdname and arguments
func Runcmd(cmdname string, args ...string) *CommandData {
	cmdpath := CheckPath(cmdname)
	var out bytes.Buffer

	cmd := exec.Command(cmdpath, args...)

	cmd.Stdout = &out

	err := cmd.Run()

	HandleError(err, "Command Failed to start => "+cmdname+" "+cmdpath)

	return &CommandData{
		Stdout:  out.Bytes(),
		Name:    cmdname,
		Process: cmd.Process,
	}

}

// A simple function to run command and return its data
// takes cmdname and arguments
func GorunCmd(ch chan<- CommandData, wg *sync.WaitGroup, cmdname string, args ...string) {
	cmdpath := CheckPath(cmdname)
	var out bytes.Buffer

	cmd := exec.Command(cmdpath, args...)

	cmd.Stdout = &out

	err := cmd.Run()

	HandleError(err, "Command Failed to start => "+cmdname+" "+cmdpath)

	datax := CommandData{
		Stdout:  out.Bytes(),
		Name:    cmdname,
		Process: cmd.Process,
	}

	ch <- datax

	defer wg.Done()

}
