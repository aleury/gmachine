package ast

import "gmachine/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
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
	}
	return ""
}

type ConstantDefinitionStatement struct {
	Token token.Token // the token.CONSTANT_DEFINITION token
	Name  *Identifier
	Value Expression
}

func (cds *ConstantDefinitionStatement) statementNode()       {}
func (cds *ConstantDefinitionStatement) TokenLiteral() string { return cds.Token.Literal }

type VariableDefinitionStatement struct {
	Token token.Token // the token.VARIABLE_DEFINITION token
	Name  *Identifier
	Value Expression
}

func (vds *VariableDefinitionStatement) statementNode()       {}
func (vds *VariableDefinitionStatement) TokenLiteral() string { return vds.Token.Literal }

type LabelDefinitionStatement struct {
	Token token.Token // the token.LABEL_DEFINITION token
}

func (lds *LabelDefinitionStatement) statementNode()       {}
func (lds *LabelDefinitionStatement) TokenLiteral() string { return lds.Token.Literal }

type OpcodeStatement struct {
	Token   token.Token // the token.OPCODE token
	Operand Expression
}

func (os *OpcodeStatement) statementNode()       {}
func (os *OpcodeStatement) TokenLiteral() string { return os.Token.Literal }

type RegisterLiteral struct {
	Token token.Token // the token.REGISTER token
}

func (rl *RegisterLiteral) expressionNode()      {}
func (rl *RegisterLiteral) TokenLiteral() string { return rl.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value uint64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

type CharacterLiteral struct {
	Token token.Token // the token.CHAR token
	Value rune
}

func (cl *CharacterLiteral) expressionNode()      {}
func (cl *CharacterLiteral) TokenLiteral() string { return cl.Token.Literal }

type StringLiteral struct {
	Token token.Token // the token.STRING token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
