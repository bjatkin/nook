package vm

import "fmt"

type Kind int64

const (
	Untyped = Kind(iota)
	Int
	Float
	Bool
	Atom
	String
	Path
	Flag
	None
)

func (r Kind) String() string {
	switch r {
	case Int:
		return "int"
	case Float:
		return "float"
	case Bool:
		return "bool"
	case Atom:
		return "atom"
	case String:
		return "string"
	case Path:
		return "path"
	case Flag:
		return "flag"
	case None:
		return "none"
	default:
		return "untyped"
	}
}

var NoneValue = Value{value: nil, kind: None}

type Value struct {
	value any
	kind  Kind
}

func (v *Value) Value() any {
	return v.value
}

func (v *Value) String() string {
	return fmt.Sprint(v.value)
}

func (v *Value) Int() int64 {
	// TODO: handle errors here?
	return v.value.(int64)
}

func (v *Value) Float() float64 {
	// TODO: handle errors here?
	return v.value.(float64)
}

func (v *Value) Kind() Kind {
	return v.kind
}
