package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	OPCODE           = "OPCODE"
	REGISTER         = "REGISTER"
	LABEL_DEFINITION = "LABEL_DEFINITION" // TODO: Add to lexer
	IDENT            = "IDENT"
	INT              = "INT"
	CHAR             = "CHAR"
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
	"DECA": OPCODE,
	"ADDA": OPCODE,
	"MOVA": OPCODE,
	"SETA": OPCODE,
	"PSHA": OPCODE,
	"POPA": OPCODE,
	"JUMP": OPCODE,
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
	return IDENT
}
