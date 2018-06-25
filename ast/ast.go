package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	// 「フィールド名を省略して埋め込まれた構造体のフィールド名が一意に定まる場合に限り、
	// 中間のフィールド名を省略してアクセスできる」を応用して、共通の性質を持たせている
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
