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
	Path
	Flag
	Atom
)

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case EOF:
		return "EOF"
	case Identifier:
		return "Identifier"
	case Comment:
		return "Comment"
	case Let:
		return "Let"
	case Exec:
		return "Exec"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Divide:
		return "Divide"
	case Multiply:
		return "Multiply"
	case GreaterThan:
		return "GreaterThan"
	case LessThan:
		return "LessThan"
	case GreaterEqual:
		return "GreaterEqual"
	case LessEqual:
		return "LessEqual"
	case Equal:
		return "Equal"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Bool:
		return "Bool"
	case String:
		return "String"
	case Path:
		return "Path"
	case Flag:
		return "Flag"
	case Atom:
		return "Atom"
	default:
		return "Uknown"
	}
}

type Token struct {
	Pos   uint
	Value string
	Kind  Kind
}
