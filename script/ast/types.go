package ast

import "github.com/bjatkin/nook/script/token"

type TypeExpr interface {
	expr()
	typeExpr()
}

type IntType struct {
	TypeExpr
	Tok token.Token
}

type FloatType struct {
	TypeExpr
	Tok token.Token
}

type BoolType struct {
	TypeExpr
	Tok token.Token
}

type AtomType struct {
	TypeExpr
	Tok token.Token
}

type StringType struct {
	TypeExpr
	Tok token.Token
}

type PathType struct {
	TypeExpr
	Tok token.Token
}

type FlagType struct {
	TypeExpr
	Tok token.Token
}

type NoneType struct {
	TypeExpr
	Tok token.Token
}

type CommandType struct {
	TypeExpr
	Tok token.Token
}

type DictType struct {
	TypeExpr
	// TODO: I want something similar to the ParamList here.
}

type TupleType struct {
	TypeExpr
	Types []TypeExpr
}

type VariadicType struct {
	TypeExpr
	Type TypeExpr
}

type FuncType struct {
	TypeExpr
	Params *ParamList
	Return TypeExpr
}

type ImplType struct {
	TypeExpr
	Funcs []FuncType
}

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
