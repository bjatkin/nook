package ast

import "github.com/bjatkin/nook/script/token"

// SExpr represents an arbitrary s-expression in the language
// Operator is an expression so it can contain things like function literals
// and identifiers bound to built-ins
type SExpr struct {
	Expr
	Operator Expr
	Operands []Expr
}

// SCommand is a command in the form '$command_name'.
// It differs from a Command expression in that it only refers to the leading
// element of the containing SExpr and not the full command expression
type SCommand struct {
	Expr
	Tok token.Token
}

// SFunc is 'fn' keyword at the beginning of an SExpr that creates a function
// literal. It differs from a Func expression in that it only refers to the leading
// element of the containing SExpr and not the full function literal
type SFunc struct {
	Expr
	Tok token.Token
}

// SLet is the 'let' keyword at the beginning of an SExpr that binds a value
// to an identifier in the parent scope.
// It differes from a Let expression in that it only refers to the leading
// element of the containing SExpr and not the full let expression
type SLet struct {
	Expr
	Tok token.Token
}

// SImpl is the 'impl' keyword at the beginning of an SExpr that binds a value
// to an identifier in the parent scope.
// It differes from an Impl expressoin in that it only refers to the leading
// element of the containing SExpr and not the full let expression
type SImpl struct {
	Expr
	Tok token.Token
}

// SSquare is the operator for s-expressions constructed using the [...] syntax.
type SSquare struct {
	Expr
	Tok token.Token
}

// SCurly is the operator for s-expressions constructed using the {...} syntax.
type SCurly struct {
	Expr
	Tok token.Token
}
