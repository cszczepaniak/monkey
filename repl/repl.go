package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/cszczepaniak/monkey/lexer"
	"github.com/cszczepaniak/monkey/parser"
)

const PROMPT = "$ "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		fmt.Fprintln(out, program.String())
	}
}
