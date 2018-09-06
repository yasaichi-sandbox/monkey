package ast

import (
	"bytes"
	"fmt"
	"github.com/yasaichi-sandbox/monkey/token"
	"strings"
)

type Node interface {
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (*BlockStatement) statementNode() {}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (*ExpressionStatement) statementNode() {}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	// NOTE: ここでnilだった場合でも`= ;`は出力されるのめちゃくちゃ微妙な気がするんだけどどうなんだろう
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (*LetStatement) statementNode() {}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (*ReturnStatement) statementNode() {}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) String() string {
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (*ArrayLiteral) expressionNode() {}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (*Boolean) expressionNode() {}

type CallExpression struct {
	Token     token.Token
	Function  Expression // FunctionLiteral or Identifier
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (*CallExpression) expressionNode() {}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fl.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (*FunctionLiteral) expressionNode() {}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) String() string {
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (*HashLiteral) expressionNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	// NOTE: 今のところ`i.Token.Literal`との違いがわかっていない
	return i.Value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (*Identifier) expressionNode() {}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (*IfExpression) expressionNode() {}

type IndexExpression struct {
	Token token.Token // '['トークン
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Index.String())
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (*IndexExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.Token // 演算子トークン、例えば「+」
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

func (*InfixExpression) expressionNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int64 // NOTE: ソースコード中の整数リテラルが表現している実際の値を格納する
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (*IntegerLiteral) expressionNode() {}

type PrefixExpression struct {
	Token    token.Token
	Operator string // "-" or "!"
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (*PrefixExpression) expressionNode() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (*StringLiteral) expressionNode() {}
