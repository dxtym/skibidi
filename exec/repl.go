package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/dxtym/skibidi/eval"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/object"
	"github.com/dxtym/skibidi/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to Ohio, %s!\n", user.Username)
	fmt.Printf("Rizz up some Skibidi yapology:\n")

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// TODO: exit function for repl
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
		io.WriteString(out, "\t" + e +" \n")
	}
	os.Exit(1)
}
