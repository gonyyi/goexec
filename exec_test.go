// (c) Gon Y. Yi 2021 <https://gonyyi.com/copyright>
// Last Update: 12/07/2021

package goexec_test

import (
	"fmt"
	"github.com/gonyyi/goexec"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a new GoExec object.
	// Note that args has to be all separated.
	// > "-s 1", "-exit 0"       // incorrect
	// > "-s", "1", "-exit", "0" // correct
	// Alternate way:
	//   c := goexec.New("./wait", "-s=4", "-exit", "0")
	c := goexec.New("./wait", "-s", "1", "-exit", "0")

	// BaseDir is where the code will be running from
	c.BaseDir = "/Users/gonyi/go/src/github.com/gonyyi/goexec/_test"

	// Async's default value is false. If it is set to true, *Exec.Run will wait until job finishes.
	// c.Async = true

	// Example of how to pipe-in data (standard input) to the process
	c.StdIn.WriteString("are you the gon?\none with yi?")

	// A function that will run when the job run successfully
	whenSucc := func(exe *goexec.Exec) {
		println()
		fmt.Printf("900  Running: %v\n", c.IsRunning())
		fmt.Printf("901  PID    : %v\n", c.PID())
		fmt.Printf("902  Err    : %v\n", c.Error())
		tmp := c.StdOut.String()
		fmt.Printf("903  StdOut : %d b\n", len(tmp))
		if len(tmp) > 0 {
			for idx, val := range strings.Split(tmp, "\n") {
				if val != "" {
					fmt.Printf("     [%d] %s\n", idx, val)
				}
			}
		}
		tmp = c.StdErr.String()
		fmt.Printf("904  StdErr : %d b\n", len(tmp))
		if len(tmp) > 0 {
			for idx, val := range strings.Split(tmp, "\n") {
				if val != "" {
					fmt.Printf("     [%d] %s\n", idx, val)
				}
			}
		}
	}

	// A function that will run when the job failed
	whenFail := func(exe *goexec.Exec) {
		fmt.Printf("700  Failed - %v\n", exe.Error())
	}

	// Run the process. This returns an error if any, however, note that if the job ran without async on,
	// the error will not return as the job is running in the background. To check the error for non async jobs,
	// Wait until the job finishes using `*Execute.IsRunning()` or `*Execute.Wait()`, then get error by
	// `*Execute.Error()`.
	c.Run(whenSucc, whenFail)

	fmt.Println("100  Running?", c.IsRunning())
	fmt.Println("101  PID?    ", c.PID())
	fmt.Println("102  Error?  ", c.Error())

	// Wait until the job running in the background process finishes.
	c.Wait()
	println("800  Done")
}
