package command

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Runner struct {
	stdout     io.Writer
	stderr     io.Writer
	workingDir string
	debug      bool
}

func Init(stdout, stderr io.Writer, dir string, debug bool) *Runner {
	return &Runner{stdout, stderr, dir, debug}
}

//Run execute the command with the given parameters
//params can't be inline, you have to send everything separated
func (r Runner) Run(command string, params ...string) error {
	cmd := exec.Command(command, params...)
	cmd.Dir = r.workingDir

	if r.debug {
		fmt.Printf("command: working directory %s \nrunner: cmd %s\n", cmd.Dir, cmd.String())
	}

	//Get readers for stdout and stderr
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err := cmd.Start()

	if err != nil {
		return fmt.Errorf("command: failed to start command:%v, err:%v", command, err)
	}

	//To synchronize before calling Wait on cmd
	var wg sync.WaitGroup
	wg.Add(2)

	var errStdout, errStderr error
	//Must be run in parallel of the command
	//to output while the command is running
	go func() {
		_, errStdout = io.Copy(r.stdout, stdoutIn)
		wg.Done()
	}()
	go func() {
		_, errStderr = io.Copy(r.stderr, stderrIn)
		wg.Done()
	}()

	wg.Wait()

	err = cmd.Wait()

	if err != nil {
		return fmt.Errorf("command: failed to run command:%v, err:%v", command, err)
	}
	if errStdout != nil || errStderr != nil {
		return fmt.Errorf("command: failed to capture stdout or stderr")
	}

	return nil
}
