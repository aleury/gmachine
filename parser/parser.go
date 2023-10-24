package parser

import (
	"errors"
	"fmt"
	"gmachine/ast"
	"gmachine/lexer"
	"gmachine/token"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidOperand error = errors.New("invalid operand")

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
	case token.OPCODE:
		return p.parseOpcodeStatement()
	case token.LABEL_DEFINITION:
		return p.parseLabelDefinitionStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLabelDefinitionStatement() ast.Statement {
	return &ast.LabelDefinitionStatement{Token: p.curToken}
}

func (p *Parser) parseOpcodeStatement() ast.Statement {
	stmt := &ast.OpcodeStatement{Token: p.curToken}

	if exprParser, ok := p.exprParsers[p.peekToken.Type]; ok {
		p.nextToken()
		stmt.Operand = exprParser()
	}

	return stmt
}

func (p *Parser) parseRegisterLiteral() ast.Expression {
	return &ast.RegisterLiteral{Token: p.curToken}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLiteral := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseUint(intLiteral.TokenLiteral(), 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("%w: %s at line %d", ErrInvalidOperand, intLiteral.TokenLiteral(), intLiteral.Token.Line))
		return nil
	}

	intLiteral.Value = value

	return intLiteral
}

func (p *Parser) parseCharacterLiteral() ast.Expression {
	charLiteral := &ast.CharacterLiteral{Token: p.curToken}

	// TODO: Handle errors
	char, size := utf8.DecodeRuneInString(strings.Trim(charLiteral.TokenLiteral(), "'"))
	if size == 0 {
		return nil
	}

	charLiteral.Value = char

	return charLiteral
}
