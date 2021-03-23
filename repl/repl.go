package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/token"
)

func StartRepl(in io.Reader, out io.Writer) {
	buf := bufio.NewScanner(in)

	for {
		fmt.Printf("[MONKEY]>>")
		scanned := buf.Scan()
		if !scanned {
			return
		}

		input := buf.Text()

		l := lexer.New(input)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
