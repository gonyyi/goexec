// (c) Gon Y. Yi 2021-2022 <https://gonyyi.com/copyright>
// Last Update: 01/13/2022

package goexec

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// Version is to identify this library.
const Version string = "GoExec v0.3.3"

var (
	ErrNoPID error = errors.New("no pid info")
)

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
	id      int          // id is optional pool ID used by bucket
	JobID   int          // JobID is optional ID for job
	BaseDir string       // BaseDir where the code should be run
	Command string       // Command is the command or file to execute
	Args    []string     // Args is command arguments
	StdIn   bytes.Buffer // StdIn if any needs to be sent to
	StdOut  bytes.Buffer // StdOut is a standard output from the run
	StdErr  bytes.Buffer // StdErr is a standard error from the run
	Async   bool         // Async will wait until the run finishes

	exitCode  int
	pid       int   // pid is for the process ID
	err       error // err if any
	isRunning bool  // isRunning stores status

	AtStart   func(*Exec)
	AtFail    func(*Exec)
	AtSuccess func(*Exec)
}

func (c *Exec) Set(baseDir, cmd string, args ...string) *Exec {
	c.BaseDir = baseDir
	c.Command = cmd
	c.Args = args
	return c
}

// ID will return Exec's ID if available
func (c *Exec) ID() int {
	return c.id
}

// Reset will reset execute
func (c *Exec) Reset() {
	if c.BaseDir != "" {
		c.BaseDir = ""
	}
	if c.Command != "" {
		c.Command = ""
	}
	c.Args = c.Args[:0]
	c.StdIn.Reset()
	c.StdOut.Reset()
	c.StdErr.Reset()
	c.Async = false
	c.pid = -1
	c.err = nil
	c.isRunning = false
}

// ExitCode returns exit code of failed job
// This won't work for Windows
func (c *Exec) ExitCode() int {
	return c.exitCode
}

// PID returns job's process ID
func (c *Exec) PID() int {
	return c.pid
}

// Error returns an error that stored when the job jobDone or jobFailed.
// To find out if the job jobDone was successful, check this value.
func (c *Exec) Error() error {
	return c.err
}

// IsRunning returns isRunning value.
// When the Run starts, this will set to true, and when jobDone this will be false.
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

// Kill will kill the currently jobRunning process
func (c *Exec) Kill() error {
	if c.PID() < 0 {
		return ErrNoPID
	}
	proc, err := os.FindProcess(c.PID())
	if err != nil {
		return err
	}
	// Kill the process
	return proc.Kill()
}

// Run will execute the command. Also, it takes a function `doneFn`,
// and execute it when the job jobDone.
// NOTE: the error this function returns is about when starting.
//   any error while jobRunning will not be returned. It can be obtained
//   from `*Exec.Error()`
func (c *Exec) Run() error {
	var err error

	cmd := exec.Command(c.Command, c.Args...)
	cmd.Dir = c.BaseDir
	cmd.Stdin = &c.StdIn
	cmd.Stdout = &c.StdOut
	cmd.Stderr = &c.StdErr

	// If starting the command jobFailed, it will end the func
	if err = cmd.Start(); err != nil {
		c.pid = -1
		c.err = err
		// FAILED
		if c.AtFail != nil {
			c.AtFail(c)
		}
		return err
	}

	c.isRunning = true
	c.pid = cmd.Process.Pid

	if c.AtStart != nil {
		c.AtStart(c)
	}

	if c.Async { // Async == true, wait until the job finishes
		c.run(cmd)
	} else { // Async == false, jobRunning as a background process
		go func() {
			c.run(cmd)
		}()
	}
	return c.Error()
}

func (c *Exec) run(cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		if msg, ok := err.(*exec.ExitError); ok { // there is error code
			c.exitCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
		}
		c.err = err
		if c.AtFail != nil {
			c.AtFail(c)
		}
	} else {
		if c.AtSuccess != nil {
			c.AtSuccess(c)
		}
	}
	c.isRunning = false
}
