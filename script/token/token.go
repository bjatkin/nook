package token

type Kind int64

const (
	Invalid = Kind(iota)
	EOF
	Identifier
	Comment

	// Keywords and Symbols
	Let
	Exec
	Plus
	Minus
	Divide
	Multiply
	GreaterThan
	LessThan
	GreaterEqual
	LessEqual
	Equal
	OpenParen
	CloseParen

	// Literals
	Int
	Float
	Bool
	String
	Atom
)

type Token struct {
	Pos   uint
	Value string
	Kind  Kind
}
