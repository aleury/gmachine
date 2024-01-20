package parser

import (
	"errors"
	"fmt"
	"gmachine/ast"
	"gmachine/lexer"
	"gmachine/token"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidSyntax error = errors.New("invalid syntax")
var ErrInvalidIntegerLiteral error = errors.New("invalid integer literal")

type expressionParserFn func() ast.Expression

type Parser struct {
	l           *lexer.Lexer
	curToken    token.Token
	peekToken   token.Token
	errors      []error
	exprParsers map[token.TokenType]expressionParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.exprParsers = make(map[token.TokenType]expressionParserFn)
	p.exprParsers[token.REGISTER] = p.parseRegisterLiteral
	p.exprParsers[token.IDENT] = p.parseIdentifier
	p.exprParsers[token.INT] = p.parseIntegerLiteral
	p.exprParsers[token.CHAR] = p.parseCharacterLiteral
	p.exprParsers[token.STRING] = p.parseStringLiteral

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.INSTRUCTION:
		return p.parseInstructionStatement()
	case token.LABEL_DEFINITION:
		return p.parseLabelDefinitionStatement()
	case token.CONSTANT_DEFINITION:
		return p.parseConstantDefinitionStatement()
	case token.VARIABLE_DEFINITION:
		return p.parseVariableDefinitionStatement()
	default:
		return nil
	}
}

func (p *Parser) parseVariableDefinitionStatement() ast.Statement {
	stmt := ast.VariableDefinitionStatement{Token: p.curToken}
	stmt.Name = p.expectOneOf(token.IDENT).(ast.Identifier)
	stmt.Value = p.expectOneOf(token.INT, token.STRING)
	return stmt
}

func (p *Parser) parseConstantDefinitionStatement() ast.Statement {
	stmt := ast.ConstantDefinitionStatement{Token: p.curToken}
	stmt.Name = p.expectOneOf(token.IDENT).(ast.Identifier)
	stmt.Value = p.expectOneOf(token.INT)
	return stmt
}

func (p *Parser) parseLabelDefinitionStatement() ast.Statement {
	return ast.LabelDefinitionStatement{Token: p.curToken}
}

func (p *Parser) expectOneOf(tokTypes ...token.TokenType) ast.Expression {
	p.nextToken()
	if !slices.Contains(tokTypes, p.curToken.Type) {
		p.errors = append(p.errors, fmt.Errorf("%w: expected one of %+v, got %s at line %d", ErrInvalidSyntax, tokTypes, p.curToken.Type, p.curToken.Line))
	}

	switch p.curToken.Type {
	case token.REGISTER:
		return p.parseRegisterLiteral()
	case token.IDENT:
		return p.parseIdentifier()
	case token.INT:
		return p.parseIntegerLiteral()
	case token.STRING:
		return p.parseStringLiteral()
	default:
		return nil
	}
}

func (p *Parser) parseInstructionStatement() ast.Statement {
	stmt := ast.InstructionStatement{Token: p.curToken}

	if stmt.TokenLiteral() == "MOVE" {
		stmt.Operand1 = p.expectOneOf(token.REGISTER, token.IDENT)
		p.expectOneOf(token.ARROW)
		stmt.Operand2 = p.expectOneOf(token.REGISTER, token.IDENT)
	}

	if exprParser, ok := p.exprParsers[p.peekToken.Type]; ok {
		p.nextToken()
		stmt.Operand1 = exprParser()
	}

	return stmt
}

func (p *Parser) parseRegisterLiteral() ast.Expression {
	return ast.RegisterLiteral{Token: p.curToken}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLiteral := ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseUint(intLiteral.TokenLiteral(), 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("%w: %s at line %d", ErrInvalidIntegerLiteral, intLiteral.TokenLiteral(), intLiteral.Token.Line))
		return nil
	}

	intLiteral.Value = value

	return intLiteral
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCharacterLiteral() ast.Expression {
	charLiteral := ast.CharacterLiteral{Token: p.curToken}

	// TODO: Handle errors
	char, size := utf8.DecodeRuneInString(strings.Trim(charLiteral.TokenLiteral(), "'"))
	if size == 0 {
		return nil
	}

	charLiteral.Value = char

	return charLiteral
}
