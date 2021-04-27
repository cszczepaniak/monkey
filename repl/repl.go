package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/cszczepaniak/monkey/evaluator"
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
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		result := evaluator.Eval(program)
		if result != nil {
			fmt.Fprintf(out, "%s\n", result.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errs []string) {
	for _, e := range errs {
		fmt.Fprintf(out, "%s\n", e)
	}
}
