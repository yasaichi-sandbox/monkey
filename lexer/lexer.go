package lexer

import (
	"github.com/yasaichi-sandbox/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置（現在の文字を指し示す）
	readPosition int  // これから読み込む位置（現在の文字の次）
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

// NOTE: 字句解析時にはこのメソッドを繰り返し呼んで行う。イテレータっぽい使われ方。
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(l.ch)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			return tok
		} else if isDigit(l.ch) {
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		}

		tok = token.Token{Type: token.ILLEGAL, Literal: ""}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	// NOTE: 終端に達した後に`readChar`を読んでも常にNUL文字を返したいならこういう実装になる
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NOTE: ASCIIのNUL文字（NUL終端文字列における文字列の終端）に対応。
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	// NOTE: for文から抜けた時のl.chは識別子/キーワードではないので、その1文字前までを取り出す
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// NOTE: そもそもなぜbyte(uint8)とrune(int32)が比較できるのかがわからなかった、、、
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}
