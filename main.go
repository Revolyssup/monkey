package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/Revolyssup/monkey/eval"
	"github.com/Revolyssup/monkey/lexer"
	"github.com/Revolyssup/monkey/obj"
	"github.com/Revolyssup/monkey/parser"
	"github.com/Revolyssup/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if len(os.Args) > 1 {
		file, err := ioutil.ReadFile(os.Args[1])
		fileact := string(file)
		if err != nil {
			panic(err)
		}
		run(fileact, os.Stdout)

		return
	}
	fmt.Printf("Welcome to monkey %s\n", user.Username)
	fmt.Printf("STARTING REPL SESSION...\n")
	repl.StartRepl(os.Stdin, os.Stdout)
}

func run(input string, out io.Writer) {

	env := obj.NewEnvironment()
	l := lexer.New(input)

	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(out, p.Errors())
	}

	evalObj := eval.Eval(program, env)
	if evalObj != nil {
		io.WriteString(out, evalObj.Inspect())
		io.WriteString(out, "\n")
	}
}

func deleteFile() {

}
