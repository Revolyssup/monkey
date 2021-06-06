package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/Revolyssup/monkey/eval"
	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/obj"
	"github.com/Revolyssup/monkey/parser"
)

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
func CloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal. Monkey says bye!")
		os.Exit(0)
	}()
}
func StartRepl(in io.Reader, out io.Writer) {
	buf := bufio.NewScanner(in)
	env := obj.NewEnvironment()
	CloseHandler()
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
			PrintParserErrors(out, p.Errors())
			continue
		}

		evalObj := eval.Eval(program, env)
		if evalObj != nil {
			io.WriteString(out, evalObj.Inspect())
			io.WriteString(out, "\n")
		}

	}
}
