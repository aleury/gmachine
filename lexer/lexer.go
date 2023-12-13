package lexer

import (
	"errors"
	"gmachine/token"
	"io"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidNumberLiteral error = errors.New("invalid number")
var ErrInvalidCharacterLiteral error = errors.New("invalid character literal, missing closing '")

type Lexer struct {
	input        string
	line         int  // current line number in input (for current char)
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	cw           int  // width of current char in bytes
}

func New(reader io.Reader) (*Lexer, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	l := &Lexer{input: string(input), line: 1}
	l.readChar()
	return l, nil
}

func (l *Lexer) NextToken() token.Token {
	for {
		l.skipWhitespace()
		switch {
		case l.ch == ';':
			l.readUntil('\n')
			continue
		case l.ch == '\'':
			char := l.readCharacter()
			return l.newToken(token.CHAR, char)
		case l.ch == 0:
			return l.newToken(token.EOF, "")
		case unicode.IsDigit(l.ch):
			value := l.readUntil('\n')
			return l.newToken(token.INT, value)
		case l.ch == '.':
			literal := l.readIdentifier()
			return l.newToken(token.LABEL_DEFINITION, literal)
		case unicode.IsLetter(l.ch):
			literal := l.readIdentifier()
			kind := token.LookupIdent(literal)
			return l.newToken(kind, literal)
		default:
			// Should we continue lexing if there is an illegal token?
			return l.newToken(token.ILLEGAL, string(l.ch))
		}
	}
}

func (l *Lexer) newToken(kind token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    kind,
		Literal: literal,
		Line:    l.line,
	}
}

func (l *Lexer) readUntil(r rune) string {
	start := l.position
	for l.ch != r && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readCharacter() string {
	start := l.position
	l.readChar()
	l.readChar()
	l.readChar()
	return l.input[start:l.position]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	if l.ch == '.' {
		l.readChar()
	}
	for unicode.IsLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch, l.cw = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.position = l.readPosition
	l.readPosition += l.cw
}
