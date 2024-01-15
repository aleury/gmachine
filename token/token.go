package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	OPCODE              = "OPCODE"
	REGISTER            = "REGISTER"
	LABEL_DEFINITION    = "LABEL_DEFINITION"
	CONSTANT_DEFINITION = "CONSTANT_DEFINITION"
	IDENT               = "IDENT"
	INT                 = "INT"
	CHAR                = "CHAR"
	STRING              = "STRING"
)

var registers = map[string]TokenType{
	"X": REGISTER,
	"Y": REGISTER,
}

var opcodes = map[string]TokenType{
	"HALT": OPCODE,
	"NOOP": OPCODE,
	"OUTA": OPCODE,
	"INCA": OPCODE,
	"INCX": OPCODE,
	"INCY": OPCODE,
	"DECA": OPCODE,
	"DECX": OPCODE,
	"DECY": OPCODE,
	"ADDA": OPCODE,
	"MULA": OPCODE,
	"MOVA": OPCODE,
	"SETA": OPCODE,
	"SETX": OPCODE,
	"SETY": OPCODE,
	"PSHA": OPCODE,
	"POPA": OPCODE,
	"JUMP": OPCODE,
	"JXNZ": OPCODE,
}

var pragmas = map[string]TokenType{
	"CONS": CONSTANT_DEFINITION,
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
