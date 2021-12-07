// (c) Gon Y. Yi 2021 <https://gonyyi.com/copyright>
// Last Update: 12/07/2021

package goexec

import (
	"bytes"
	"github.com/gonyyi/gosl"
	"os/exec"
)

// Version is to identify this library.
const Version gosl.Ver = "GoExec v0.1.0-1"

// New creates Exec and return its pointer
// Exec can be recycled thru `*Exec.Reset()`
func New(command string, args ...string) *Exec {
	return &Exec{
		Command: command,
		Args:    args,
	}
}

// Exec runs command with optional arguments and piped standard input.
type Exec struct {
	BaseDir   string       // BaseDir where the code should be run
	Command   string       // Command is the command or file to execute
	Args      []string     // Args is command arguments
	StdIn     bytes.Buffer // StdIn if any needs to be sent to
	StdOut    bytes.Buffer // StdOut is a standard output from the run
	StdErr    bytes.Buffer // StdErr is a standard error from the run
	Async     bool         // Async will wait until the run finishes
	exitCode  int
	pid       int   // pid is for the process ID
	err       error // err if any
	isRunning bool  // isRunning stores status
}

// Reset will reset execute
func (c *Exec) Reset() {
	c.BaseDir = ""
	c.Command = ""
	c.Args = c.Args[:0]
	c.StdIn.Reset()
	c.StdOut.Reset()
	c.StdErr.Reset()
	c.Async = false
	c.pid = -1
	c.err = nil
	c.isRunning = false
}

// PID returns job's process ID
func (c *Exec) PID() int {
	return c.pid
}

// Error returns an error that stored when the job finished or failed.
// To find out if the job finished was successful, check this value.
func (c *Exec) Error() error {
	return c.err
}

// IsRunning returns isRunning value.
// When the Run starts, this will set to true, and when finished this will be false.
func (c *Exec) IsRunning() bool {
	return c.isRunning
}

// Wait will wait until the non-Async job to finish
func (c *Exec) Wait() {
	for {
		if c.isRunning == false {
			break
		}
	}
}

// Run will execute the command. Also, it takes a function `doneFn`,
// and execute it when the job finished.
func (c *Exec) Run(whenSucc, whenFail func(*Exec)) error {
	var err error

	cmd := exec.Command(c.Command, c.Args...)
	cmd.Dir = c.BaseDir
	cmd.Stdin = &c.StdIn
	cmd.Stdout = &c.StdOut
	cmd.Stderr = &c.StdErr

	// If starting the command failed, it will end the func
	if err = cmd.Start(); err != nil {
		c.pid = -1
		c.err = err
		return err
	}

	// Since there isn't an error starting the command,
	// set the isRunning = true,
	// set the pid
	c.isRunning = true
	c.pid = cmd.Process.Pid

	do := func() {
		if err = cmd.Wait(); err != nil {
			c.err = err
			if whenFail != nil {
				whenFail(c)
			}
		} else {
			if whenSucc != nil {
				whenSucc(c)
			}
		}
		c.isRunning = false
	}
	// If Async == true, wait until the process is finished
	// Otherwise run it in background
	if c.Async { // Async == true, wait until the job finishes
		do()
	} else { // Async == false, running as a background process
		go func() {
			do()
		}()
	}
	return c.Error()
}
