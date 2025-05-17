package ast

import "github.com/bjatkin/nook/script/token"

type TypeExpr interface {
	expr()
	typeExpr()
}

// IntType represents the `int` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type IntType struct {
	TypeExpr
	Tok token.Token
}

// FloatType represents the `float` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type FloatType struct {
	TypeExpr
	Tok token.Token
}

// BoolType represents the `bool` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type BoolType struct {
	TypeExpr
	Tok token.Token
}

// AtomType represents the `atom` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type AtomType struct {
	TypeExpr
	Tok token.Token
}

// StringType represents the `str` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type StringType struct {
	TypeExpr
	Tok token.Token
}

// PathType represents the `path` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type PathType struct {
	TypeExpr
	Tok token.Token
}

// FlagType represents the `flag` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type FlagType struct {
	TypeExpr
	Tok token.Token
}

// NoneType represents the `none` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type NoneType struct {
	TypeExpr
	Tok token.Token
}

// CommandType represents the `cmd` keyword in a type expression in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type CommandType struct {
	TypeExpr
	Tok token.Token
}

// DictType represents a dictionary type in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type DictType struct {
	TypeExpr
	// TODO: I want something similar to the ParamList here.
}

// TupleType represents a tuple type in NookScript.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type TupleType struct {
	TypeExpr
	Types []TypeExpr
}

// VariadicType represents a variadic type in NookScript fuction paramater list.
// It is only allowed in the final position of the paramater list.
// Types can always be omitted and then infered in NookScript, in which case
// this node will not be added until the normalizer or checker phases
type VariadicType struct {
	TypeExpr
	Type TypeExpr
}

// FuncType represents a functions type in NookScript.
// The function must, at a minimum, include a list of paramaters names.
// If paramater types are not provided they will be infered based on their
// usage inside the functions body.
type FuncType struct {
	TypeExpr
	Params *ParamList
	Return TypeExpr
}

// ImplType represents an implementation type that wraps one or more function overrides.
// This allows a function to act polymorphicly based on the input it recieves.
type ImplType struct {
	TypeExpr
	Funcs []FuncType
}

// TraitType represents a trait in NookScript which allows types to be constrained by
// Behavior, rather than a simple type expression.
type TraitType struct {
	TypeExpr
	// TODO: figure out how traits are going to be represented
	// for not all traits will just be considered to be the empyt
	// trait or the 'any' type
	//
	// Traits should include
	// * Behavior
	//     * Function calls (including joint types)
	//     * Property access
	//     * Slice access
	// * Type Constraints
}
