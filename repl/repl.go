package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nomad-software/script/lexer"
)

const (
	prompt = ">>> "
)

// Start the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		tokens := lexer.New(line)

		for token := range tokens {
			fmt.Fprintf(out, "%s\n", token)
		}
	}
}
