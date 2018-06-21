package lexer_test

import (
	"github.com/yasaichi-sandbox/monkey/lexer"
	"github.com/yasaichi-sandbox/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{ // NOTE: フィールド名を明示しないで初期化する場合、定義順に値を並べる
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	// NOTE: mapを使ってsliceを作り、これが期待通りか確認するとやりたいところだがしょうがない
	l := lexer.New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf( // NOTE: `Fatalf`を使ってassertionが失敗したときに以降のテストを実行しないようにする
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
