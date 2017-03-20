package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/parser"
)

// Start the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(">>> ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.Parse()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
