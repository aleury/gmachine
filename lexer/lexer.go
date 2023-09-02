package lexer

import (
	"gmachine/token"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	cw           int  // width of current char in bytes
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	for {
		l.skipWhitespace()
		switch l.ch {
		case ';':
			l.readUntil('\n')
			continue
		case '\'':
			return token.Token{Type: token.CHAR, Literal: l.readCharacter()}
		case 0:
			return token.Token{Type: token.EOF, Literal: ""}
		default:
			var tok token.Token
			if unicode.IsDigit(l.ch) {
				tok.Type = token.INT
				tok.Literal = l.readInt()
			} else if unicode.IsLetter(l.ch) || (l.ch == '.' && unicode.IsLetter(l.peek())) {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
			} else {
				tok.Type = token.ILLEGAL
				tok.Literal = string(l.ch)
			}
			return tok
		}
	}
}

func (l *Lexer) readUntil(r rune) string {
	start := l.position
	for l.ch != r {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) peek() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	nextChar, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return nextChar
}

func (l *Lexer) readCharacter() string {
	start := l.position
	l.readChar()
	if l.peek() != '\'' {
		// TODO(adam): Handle this better by returning error?
		panic("invalid character literal")
	}
	l.readChar()
	l.readChar()
	return l.input[start:l.position]
}

func (l *Lexer) readInt() string {
	start := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
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
