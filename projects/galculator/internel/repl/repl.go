package repl

import (
	"bufio"
	"fmt"
	"galculator/internel/compute"
	"os"
	"runtime"
)

// REPL stands for read, evaluate, print, loop.
// It's an interactive console so to speak.
func REPL() error {
	buf := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			return err
		}

		trim := 1
		if runtime.GOOS == "windows" {
			trim = 2
		}

		if len(sentence) <= trim {
			continue
		}
		fmt.Println(compute.Compute(string(sentence[:len(sentence)-trim])))
	}
}
