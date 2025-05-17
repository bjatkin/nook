package normalizer

import (
	"fmt"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

// Normalizer normalizes the initial ast into a semantically correct ast
type Normalizer struct {
	Errors []error
}

// addError adds a new error to the Normalizer struct
func (n *Normalizer) addError(err error) {
	n.Errors = append(n.Errors, err)
}

// Normalize takes an expression and normalizes it into a semantically correct ast node
func (n *Normalizer) Normalize(expr ast.Expr) ast.Expr {
	switch expr := expr.(type) {
	case *ast.SExpr:
		normalized, err := n.normalizeSExpr(expr.Operator, expr.Operands...)
		if err != nil {
			n.addError(err)
			return nil
		}

		return normalized
	default:
		return expr
	}
}

// normalizeSExpr converts s-expression into a more specific ast node
func (n *Normalizer) normalizeSExpr(operator ast.Expr, operands ...ast.Expr) (ast.Expr, error) {
	switch operator := operator.(type) {
	case *ast.SCommand:
		// convert from $git -> git
		name := operator.Tok.Value[1:]

		normArgs := []ast.Expr{}
		for _, op := range operands {
			expr := n.Normalize(op)
			normArgs = append(normArgs, expr)
		}

		return &ast.Command{
			Tok:  operator.Tok,
			Name: name,
			Args: normArgs,
		}, nil
	case *ast.SLet:
		if len(operands) != 2 {
			return nil, fmt.Errorf("let expression takes 2 operands (let [identifier] [value])")
		}

		identifier, ok := operands[0].(*ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("firt operand to 'let' must be an identifier but got '%T'", operands[0])
		}

		return &ast.Let{
			Tok:        operator.Tok,
			Identifier: identifier,
			Value:      n.Normalize(operands[1]),
		}, nil
	case *ast.SFunc:
		// support functions in the form (fn [params] type [body])
		// as well as (fn [params] [body]) where the return type is infered
		switch len(operands) {
		case 2:
			return n.normalizeUntypedFunc(operator.Tok, operands...)
		case 3:
			return n.normalizeTypedFunc(operator.Tok, operands...)
		default:
			return nil, fmt.Errorf("fn expression takes either 3 or 4 operands (fn [params] <return type> [body])")
		}
	case *ast.Identifier:
		// assume this is a function call, this will be validated in the type checker since we need to evaluate
		// identifier types before we can know the identifiers type for certian
		normArgs := []ast.Expr{}
		for _, op := range operands {
			expr := n.Normalize(op)
			normArgs = append(normArgs, expr)
		}

		return &ast.Call{
			Func: operator,
			Args: normArgs,
		}, nil
	case *ast.SExpr:
		normalizedOp, err := n.normalizeSExpr(operator.Operator, operator.Operands...)
		if err != nil {
			return nil, err
		}

		normArgs := []ast.Expr{}
		for _, op := range operands {
			expr := n.Normalize(op)
			normArgs = append(normArgs, expr)
		}

		// Functions literals can be called directly if they are the s-expression operator
		if fn, ok := normalizedOp.(*ast.Func); ok {
			return &ast.Call{
				Func: fn,
				Args: normArgs,
			}, nil
		}

		// Functions calls can be the operator as long as they return a valid function value
		// Type checking will happen later, for now assume it's a valid call
		if call, ok := normalizedOp.(*ast.Call); ok {
			return &ast.Call{
				Func: call,
				Args: normArgs,
			}, nil
		}

		return nil, fmt.Errorf("invalid s-expression operator '%v'", operator)
	default:
		return nil, fmt.Errorf("unknown s-expression operator '%v'", operator)
	}
}

// normalizeUntypedFunc normalizes s-expressions in the form (fn [params] (body)) into a function literal
func (n *Normalizer) normalizeUntypedFunc(fn token.Token, operands ...ast.Expr) (*ast.Func, error) {
	if len(operands) != 2 {
		return nil, fmt.Errorf("expected expression in the form (fn [<params>] (body))")
	}

	params, ok := operands[0].(ast.SExpr)
	if !ok {
		return nil, fmt.Errorf("first argument to a function definition must be a parameter list not '%v'", operands[0])
	}
	_, ok = params.Operator.(ast.SSquare)
	if !ok {
		return nil, fmt.Errorf("first argument to a function definition must be a parameter list not '%v'", operands[0])
	}

	paramList, err := normalizeParamList(params.Operands)
	if err != nil {
		return nil, fmt.Errorf("invalid paramater list for function %w", err)
	}

	body := operands[1]

	return &ast.Func{
		Tok: fn,
		Type: &ast.FuncType{
			Params: paramList,
			Return: &ast.TraitType{},
		},
		Body: body,
	}, nil
}

// normalize s-expression in the form (fn [params] type (body)) into a function literal
func (n *Normalizer) normalizeTypedFunc(fn token.Token, operands ...ast.Expr) (*ast.Func, error) {
	if len(operands) != 3 {
		return nil, fmt.Errorf("expected expression in the form (fn [<param>] type (body))")
	}

	params, ok := operands[0].(ast.SExpr)
	if !ok {
		return nil, fmt.Errorf("first argument to a function definition must be a parameter list not '%v'", operands[0])
	}
	_, ok = params.Operator.(ast.SSquare)
	if !ok {
		return nil, fmt.Errorf("first argument to a function definition must be a parameter list not '%v'", operands[0])
	}

	paramList, err := normalizeParamList(params.Operands)
	if err != nil {
		return nil, fmt.Errorf("invalid paramater list for function %w", err)
	}

	returnType, ok := operands[1].(ast.TypeExpr)
	if !ok {
		return nil, fmt.Errorf("second argument to a function definition must be a return type '%v'", operands[1])
	}

	body := operands[2]

	return &ast.Func{
		Tok: fn,
		Type: &ast.FuncType{
			Params: paramList,
			Return: returnType,
		},
		Body: body,
	}, nil
}

// normalizeParamList normalizes a paramater list in the form [ident type ...] or [ident ...]
func normalizeParamList(exprs []ast.Expr) (*ast.ParamList, error) {
	if len(exprs) == 0 {
		return &ast.ParamList{}, nil
	}

	identifiers := []*ast.Identifier{}
	types := []ast.TypeExpr{}
	for _, expr := range exprs {
		if param, ok := expr.(*ast.Identifier); ok {
			identifiers = append(identifiers, param)
			continue
		}

		// this allows for the syntatic shorthand [a, b, c int] where 'a' 'b' and 'c'
		// are all typed as integers
		if typeExpr, ok := expr.(ast.TypeExpr); ok {
			for len(types) < len(identifiers) {
				types = append(types, typeExpr)
			}
			continue
		}

		return nil, fmt.Errorf("invalid expression in param list '%v'", expr)
	}

	// if there are no types in the paramater list assume all the types must be infered
	if len(types) == 0 {
		for len(types) < len(identifiers) {
			types = append(types, &ast.TraitType{})
		}
	}

	paramList := &ast.ParamList{}
	for i := range identifiers {
		param := ast.Param{
			Identifier: identifiers[i],
			Type:       types[i],
		}
		paramList.Params = append(paramList.Params, param)
	}

	return paramList, nil
}
