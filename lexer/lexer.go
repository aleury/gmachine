package lexer

import (
	"errors"
	"gmachine/token"
	"io"
	"unicode"
)

var ErrInvalidNumberLiteral error = errors.New("invalid number")
var ErrInvalidCharacterLiteral error = errors.New("invalid character literal, missing closing '")

type Lexer struct {
	input         []rune
	line          int  // current line number in input (for current char)
	position      int  // current position in input (points to current char)
	nextRuneIndex int  // current reading position in input (after current char)
	currentRune   rune // current char under examination
}

func New(reader io.Reader) (*Lexer, error) {
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	l := &Lexer{input: []rune(string(input)), line: 1}
	l.readChar()
	return l, nil
}

func (l *Lexer) NextToken() token.Token {
	for {
		l.skipWhitespace()
		switch {
		case l.currentRune == ';':
			l.readUntil('\n')
			continue
		case l.currentRune == '\'':
			char := l.readCharacter()
			return l.newToken(token.CHAR, char)
		case l.currentRune == '"':
			str := l.readString()
			return l.newToken(token.STRING, str)
		case l.currentRune == '-':
			if l.peekChar() == '>' {
				l.readChar()
				l.readChar()
				return l.newToken(token.ARROW, "->")
			}
			return l.newToken(token.ILLEGAL, string(l.currentRune))
		case l.currentRune == 0:
			return l.newToken(token.EOF, "")
		case unicode.IsDigit(l.currentRune):
			value := l.readUntil('\n')
			return l.newToken(token.INT, value)
		case l.currentRune == '.':
			literal := l.readIdentifier()
			return l.newToken(token.LABEL_DEFINITION, literal)
		case unicode.IsLetter(l.currentRune):
			literal := l.readIdentifier()
			kind := token.LookupIdent(literal)
			return l.newToken(kind, literal)
		default:
			// Should we continue lexing if there is an illegal token?
			return l.newToken(token.ILLEGAL, string(l.currentRune))
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
	for l.currentRune != r && l.currentRune != 0 {
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.currentRune == '"' || l.currentRune == 0 {
			break
		}
	}
	// consume closing quote
	l.readChar()
	return string(l.input[position : l.position-1])
}

func (l *Lexer) readCharacter() string {
	start := l.position
	l.readChar()
	l.readChar()
	l.readChar()
	return string(l.input[start:l.position])
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	if l.currentRune == '.' {
		l.readChar()
	}
	for unicode.IsLetter(l.currentRune) {
		l.readChar()
	}
	return string(l.input[start:l.position])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.currentRune) {
		if l.currentRune == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	l.currentRune = l.peekChar()
	l.position = l.nextRuneIndex
	l.nextRuneIndex = l.position + 1
}

func (l *Lexer) peekChar() rune {
	if l.nextRuneIndex >= len(l.input) {
		return 0
	}
	return l.input[l.nextRuneIndex]
}
