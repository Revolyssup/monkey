package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/parser"
)

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
		for _, stmt := range program.Statements {
			fmt.Print(stmt.TokenLiteral())
		}
		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	fmt.Printf("%+v\n", tok)
		// }
	}
}
