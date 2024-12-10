package repl

import (
	"bufio"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment(nil)
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "quit" {
			return
		}
		lexer := lexer.New(line)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}
		//io.WriteString(out, program.String())

		eval := evaluator.Eval(program, env)
		if eval == evaluator.NULL {
			continue
		}
		if eval != nil {
			io.WriteString(out, eval.Inspect())
			io.WriteString(out, "\n")
		} else {
			io.WriteString(out, "sorry,can't realize this now\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
