package repl

import (
	"bufio"
	"fmt"
	"github.com/yasaichi-sandbox/monkey/lexer"
	"github.com/yasaichi-sandbox/monkey/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Print(PROMPT)

	for scanner.Scan() {
		line := scanner.Text()
		l := lexer.New(line)

		// NOTE: 最初わからなかったけど、`i := 0; i < n; i++`などと同じ構造をしている
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}

		fmt.Print(PROMPT)
	}
}
