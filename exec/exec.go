package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dxtym/skibidi/eval"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/object"
	"github.com/dxtym/skibidi/parser"
)

const (
	EXT    = ".skbd"
	PROMPT = ">> "
)

func Run(in io.Reader, out io.Writer, args []string) {
	env := object.NewEnvironment()

	if len(args) > 1 {
		runFile(out, env, args[1])
	} else {
		runRepl(in, out, env)
	}
}

func runFile(out io.Writer, env *object.Environment, file string) {
	if filepath.Ext(file) != EXT {
		io.WriteString(out, "red flag")
		os.Exit(1)
	}

	text, err := os.ReadFile(file)
	if err != nil {
		io.WriteString(out, "skill issue")
		os.Exit(1)
	}

	parseProgram(out, env, string(text))
}

func runRepl(in io.Reader, out io.Writer, env *object.Environment) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to Ohio, %s!\n", user.Username)
	fmt.Printf("Rizz up some Skibidi:\n")

	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		parseProgram(out, env, scanner.Text())
	}
}

func parseProgram(out io.Writer, env *object.Environment, text string) {
	l := lexer.NewLexer(text)
	p := parser.NewParser(l)

	program := p.Parse()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			io.WriteString(out, "\t"+e+"\t")
		}
		os.Exit(1)
	}

	evaled := eval.Eval(program, env)
	if evaled != nil {
		io.WriteString(out, evaled.Inspect())
		io.WriteString(out, "\n")
	}
}
