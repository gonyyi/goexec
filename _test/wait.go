package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/gonyyi/gosl"
	"log"
	"os"
	"time"
)

func main() {
	var sleep = 3
	var exitCode = 0

	flag.IntVar(&sleep, "s", sleep, "sleep in second")
	flag.IntVar(&exitCode, "exit", exitCode, "exit code")
	flag.Parse()

	log := gosl.NewLogger(os.Stdout).SetPrefix("wait> ")

	log.String("start")
	log.KeyInt("set sleep", sleep)
	log.KeyInt("set exit code", exitCode)

	log.KeyBool("os.Stdin == nil", os.Stdin == nil)

	b := ReadPipe()
	log.KeyInt("stdin read", len(b))

	if len(b) > 0 {
		sb := bytes.Split(b, []byte("\n"))
		for idx, v := range sb {
			if len(v) > 0 {
				log.String("  [" + gosl.Itoa(idx) + "] " + string(v))
			}
		}
	}

	log.KeyInt("go sleep for", sleep)
	time.Sleep(time.Second * time.Duration(sleep))
	log.String("end")
	os.Stderr.WriteString("wait> test for stderr\n")
	os.Stdout.WriteString("wait> test for stdout\n")
	os.Exit(exitCode)
}

func ReadPipe() []byte {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var stdin []byte
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
			stdin = append(stdin, '\n')
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return stdin
	}
	return nil
}
