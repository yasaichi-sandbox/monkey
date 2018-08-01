package repl

import (
	"bufio"
	"fmt"
	"github.com/yasaichi-sandbox/monkey/lexer"
	"github.com/yasaichi-sandbox/monkey/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		fmt.Fprintln(out, program.String())
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, MONKEY_FACE)
	fmt.Fprintln(out, "Woops! We ran into some monkey business here!")
	fmt.Fprintln(out, "parser errors:")

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
