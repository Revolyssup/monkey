package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Revolyssup/monkey/eval"
	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/parser"
)

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func StartRepl(in io.Reader, out io.Writer) {
	buf := bufio.NewScanner(in)

	for {
		fmt.Printf("\n[MONKEY]>>")
		scanned := buf.Scan()
		if !scanned {
			return
		}

		input := buf.Text()

		l := lexer.New(input)

		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evalObj := eval.Eval(program)
		if evalObj != nil {
			io.WriteString(out, evalObj.Inspect())
			io.WriteString(out, "\n")
		}

	}
}
