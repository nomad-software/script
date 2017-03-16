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
		lexer := lexer.New(line)

		for tok := range lexer.Tokens {
			fmt.Fprintf(out, "%s\n", tok)
		}
	}
}
