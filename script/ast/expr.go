package ast

import (
	"github.com/bjatkin/nook/script/token"
)

// Expr is the interface that types of all expressions in nook script
type Expr interface {
	expr()
}

// Int is an integer literal (e.g. 42)
type Int struct {
	Expr
	Tok   token.Token
	Value int64
}

// Float is a floating point literal (e.g. 3.14)
type Float struct {
	Expr
	Tok   token.Token
	Value float64
}

// String is a string literal (e.g. "hello there")
type String struct {
	Expr
	Tok   token.Token
	Value string
}

// Atom is an atom literal (e.g. 'ok)
type Atom struct {
	Expr
	Tok   token.Token
	Value string
}

// Bool is a boolean literal (e.g. true, false)
type Bool struct {
	Expr
	Tok   token.Token
	Value bool
}

// Path is a path literal (e.g. ./root/dir)
type Path struct {
	Expr
	Tok   token.Token
	Value string
}

// Flag is a flag literal (e.g. --version)
type Flag struct {
	Expr
	Tok   token.Token
	Value string
}

// Func is a function literal (e.g. (fn (a int, b int) int (+ a b))
type Func struct {
	Expr
	Tok  token.Token
	Type *FuncType
	Body Expr
}

// Identifier is an identifier in the language (e.g. user_name)
type Identifier struct {
	Expr
	Tok  token.Token
	Name string
}

// Let is a full let expression in the language (e.g. (let name "Jill")).
// It is different from SLet as it encompases the full expression and not just the
// 'let' keyword at the begining of the SExpr.
type Let struct {
	Expr
	Tok        token.Token
	Identifier *Identifier
	Value      Expr
}

// Builtin represents a builtin function. It is not representable in the language but
// it can be called in a Call expression.
type Builtin struct {
	Expr
	Fn func(args ...any) (any, error)
}

// Impl is a full imple expression in the language (e.g. (impl add (fn [a b] (+ a b))))
// It is different from SImpl as it encompases the full expression and not just the
// 'impl' keyword at the begining of the SImpl.
type Impl struct {
	Expr
	Tok        token.Token
	Identifier *Identifier
	Func       *Func
}

// Command is a call to a command (e.g. ($git 'status))
// It is different from SCommand as it encompases the full expression and not just
// the Command keyword at the begining of the SExpr.
type Command struct {
	Expr
	Tok  token.Token
	Name string
	Args []Expr
	// TODO: redirects?
}

// Call is a function call (e.g. (print "hello there"))
type Call struct {
	Expr
	Func Expr
	Args []Expr
}

// Param is a typed parameter used in function literals
// It is not an expression as it can only appear inside a ParamList
type Param struct {
	Identifier *Identifier
	Type       TypeExpr
}

// ParamList is a list of parameters that will be bound on a call to a function.
// It is an expression because it must appear as an element of an SExpr when a
// function literal is declared.
// However, it semantic validitiy is limited compared to most expressions and the
// type checker will ensure it's used properly
type ParamList struct {
	Expr
	Params []Param
}
