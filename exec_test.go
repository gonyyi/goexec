// (c) Gon Y. Yi 2021-2022 <https://gonyyi.com/copyright>
// Last Update: 01/13/2022


package goexec_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gonyyi/goexec"
)

func printMultiple(s string) {
	if len(s) > 0 {
		for idx, val := range strings.Split(s, "\n") {
			if val != "" {
				fmt.Printf("     [%d] %s\n", idx, val)
			}
		}
	}
}

func TestNew(t *testing.T) {
	// Create a new GoExec object.
	c := goexec.New("sleep", "15")
	// c.BaseDir = "/Users/gonyi/go/src/github.com/gonyyi/goexec" // optional

	// Async's default value is false. If it is set to true, *Exec.AddJob will wait until job finishes.
	// That means PID won't be useful.
	// c.Async = true

	// Example of how to pipe-in data (standard input) to the process
	// c.StdIn.WriteString("are you the gon?\none with yi?")

	atStart := func(exe *goexec.Exec) {
		fmt.Printf("Started => JobID=%d; PID=%d\n", c.JobID, c.PID())
	}

	// A function that will run when the job run successfully
	atSucc := func(exe *goexec.Exec) {
		println()

		fmt.Printf("SUCC 900  JobRunning: %v\n", c.IsRunning())
		fmt.Printf("SUCC 901  PID    : %v\n", c.PID())
		fmt.Printf("SUCC 902  Err    : %v\n", c.Error())
		tmp := c.StdOut.String() // shows StdOut
		fmt.Printf("SUCC 903  StdOut : %d b\n", len(tmp))
		printMultiple(tmp)

		tmp = c.StdErr.String() // shows StdErr
		fmt.Printf("SUCC 904  StdErr : %d b\n", len(tmp))
		printMultiple(tmp)
	}

	// A function that will run when the job jobFailed
	atFail := func(exe *goexec.Exec) {
		fmt.Printf("FAIL 700  JobFailed - %v\n", exe.Error())

		tmp := c.StdOut.String() // shows StdOut
		fmt.Printf("FAIL 703  JobFailed StdOut : %d b\n", len(tmp))
		printMultiple(tmp)

		tmp = c.StdErr.String() // shows StdErr
		fmt.Printf("FAIL 704  JobFailed StdErr : %d b\n", len(tmp))
		printMultiple(tmp)
	}

	// AddJob the process. This returns an error if any, however, note that if the job ran without async on,
	// the error will not return as the job is jobRunning in the background. To check the error for non async jobs,
	// Wait until the job finishes using `*Execute.Closed()` or `*Execute.Wait()`, then get error by
	// `*Execute.Error()`.
	c.AtStart, c.AtFail, c.AtSuccess = atStart, atFail, atSucc

	// Start
	if err := c.Run(); err != nil {
		println(err.Error())
	}

	fmt.Println("100  JobRunning?", c.IsRunning())
	fmt.Println("101  PID?    ", c.PID())
	fmt.Println("102  Error?  ", c.Error())

	// Wait until the job jobRunning in the background process finishes.
	go func() {
		time.Sleep(time.Second * 5)
		println("Kill after 5 sec")
		if err := c.Kill(); err != nil {
			println(err.Error())
		}
	}()
	c.Wait()
	println("800  Finished")
}
