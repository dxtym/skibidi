package exec

import (
	"io"
	"os"
	"path/filepath"

	"github.com/dxtym/skibidi/eval"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/object"
	"github.com/dxtym/skibidi/parser"
)

var ext = ".skbd"

func Run(in io.Reader, out io.Writer, file string) {
	if filepath.Ext(file) != ext {
		printFileErrors(out, "bruh: file not supported")
	}

	content, err := os.ReadFile(file)
	if err != nil {
		printFileErrors(out, "bruh: cannot read file")
	}

	env := object.NewEnvironment()
	l := lexer.NewLexer(string(content))
	p := parser.NewParser(l)

	program := p.Parse()
	if len(p.Errors()) != 0 {
		printParseErrors(out, p.Errors())
	}

	evaled := eval.Eval(program, env)
	if evaled != nil {
		io.WriteString(out, evaled.Inspect())
		io.WriteString(out, "\n")
	}
}

func printFileErrors(out io.Writer, err string) {
	io.WriteString(out, "\t"+err+"\n")
	os.Exit(1)
}
