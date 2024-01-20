package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INSTRUCTION         = "INSTRUCTION"
	REGISTER            = "REGISTER"
	LABEL_DEFINITION    = "LABEL_DEFINITION"
	CONSTANT_DEFINITION = "CONSTANT_DEFINITION"
	VARIABLE_DEFINITION = "VARIABLE_DEFINITION"
	IDENT               = "IDENT"
	INT                 = "INT"
	CHAR                = "CHAR"
	STRING              = "STRING"
	ARROW               = "ARROW"
	ASTERISK            = "ASTERISK"
)

var registers = map[string]TokenType{
	"A": REGISTER,
	"X": REGISTER,
	"Y": REGISTER,
}

var opcodes = map[string]TokenType{
	"HALT": INSTRUCTION,
	"NOOP": INSTRUCTION,
	"MOVE": INSTRUCTION,
	"OUTA": INSTRUCTION,
	"INCA": INSTRUCTION,
	"INCX": INSTRUCTION,
	"INCY": INSTRUCTION,
	"DECA": INSTRUCTION,
	"DECX": INSTRUCTION,
	"DECY": INSTRUCTION,
	"ADDA": INSTRUCTION,
	"MULA": INSTRUCTION,
	"SETA": INSTRUCTION,
	"SETX": INSTRUCTION,
	"SETY": INSTRUCTION,
	"PSHA": INSTRUCTION,
	"POPA": INSTRUCTION,
	"JUMP": INSTRUCTION,
	"JXNZ": INSTRUCTION,
}

var pragmas = map[string]TokenType{
	"CONS": CONSTANT_DEFINITION,
	"VARB": VARIABLE_DEFINITION,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string // Possibily rename to Value
	Line    int
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := opcodes[ident]; ok {
		return tokType
	}
	if tokType, ok := registers[ident]; ok {
		return tokType
	}
	if tokType, ok := pragmas[ident]; ok {
		return tokType
	}
	return IDENT
}
