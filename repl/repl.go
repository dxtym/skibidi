package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dxtym/maymun/eval"
	"github.com/dxtym/maymun/lexer"
	"github.com/dxtym/maymun/object"
	"github.com/dxtym/maymun/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaled := eval.Eval(program, env)
		if evaled != nil {
			io.WriteString(out, evaled.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, err []string) {
	for _, e := range err {
		io.WriteString(out, "\t"+e+"\n")
	}
}
