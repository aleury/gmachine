package lexer

import (
	"errors"
	"fmt"
	"gmachine/token"
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

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (token.Token, error) {
	for {
		l.skipWhitespace()
		switch {
		case l.ch == ';':
			l.readUntil('\n')
			continue
		case l.ch == '\'':
			char, err := l.readCharacter()
			if err != nil {
				return token.Token{}, err
			}
			return l.newToken(token.CHAR, char), nil
		case l.ch == 0:
			return l.newToken(token.EOF, ""), nil
		case unicode.IsDigit(l.ch):
			value, err := l.readInt()
			if err != nil {
				return token.Token{}, err
			}
			return l.newToken(token.INT, value), nil
		case l.ch == '.':
			if !unicode.IsLetter(l.peekChar()) {
				return token.Token{}, fmt.Errorf("invalid label definition")
			}
			literal := l.readIdentifier()
			return l.newToken(token.LABEL_DEFINITION, literal), nil
		case unicode.IsLetter(l.ch):
			literal := l.readIdentifier()
			kind := token.LookupIdent(literal)
			return l.newToken(kind, literal), nil
		default:
			// Should we continue lexing if there is an illegal token?
			return l.newToken(token.ILLEGAL, string(l.ch)), nil
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

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	nextChar, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return nextChar
}

func (l *Lexer) readCharacter() (string, error) {
	start := l.position
	l.readChar()
	if l.peekChar() != '\'' {
		return "", fmt.Errorf("%w: %s at line %d", ErrInvalidCharacterLiteral, l.input[start:l.readPosition], l.line)
	}
	l.readChar()
	l.readChar()
	return l.input[start:l.position], nil
}

func (l *Lexer) readInt() (string, error) {
	start := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	if unicode.IsLetter(l.ch) {
		return "", fmt.Errorf("%w: %s at line %d", ErrInvalidNumberLiteral, l.input[start:l.readPosition], l.line)
	}
	return l.input[start:l.position], nil
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
