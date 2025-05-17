package token

type Kind int64

const (
	Invalid = Kind(iota)
	EOF
	Identifier
	Comment
	Whitespace

	// Keywords and Symbols
	Let
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
	OpenCurly
	CloseCurly
	OpenSquare
	CloseSquare

	// Type Keywords
	IntType
	FloatType
	BoolType
	StringType
	PathType
	FlagType
	AtomType
	CommandType
	DictType
	TupleType
	SliceType
	ArrayType
	NoneType

	// Literals
	Int
	Float
	Bool
	String
	Path
	Flag
	Atom
	Command
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
	case OpenCurly:
		return "OpenCurly"
	case CloseCurly:
		return "CloseCurly"
	case OpenSquare:
		return "OpenSquare"
	case CloseSquare:
		return "CloseSquare"
	case IntType:
		return "IntType"
	case FloatType:
		return "FloatType"
	case BoolType:
		return "BoolType"
	case StringType:
		return "StringType"
	case PathType:
		return "PathType"
	case FlagType:
		return "FlagType"
	case AtomType:
		return "AtomType"
	case CommandType:
		return "CommandType"
	case DictType:
		return "DictType"
	case TupleType:
		return "TupleType"
	case SliceType:
		return "SliceType"
	case ArrayType:
		return "ArrayType"
	case NoneType:
		return "NoneType"
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
	case Command:
		return "Command"
	default:
		return "Uknown"
	}
}

type Token struct {
	Pos   uint
	Value string
	Kind  Kind
}
