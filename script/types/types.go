package types

import (
	"github.com/bjatkin/nook/script/ast"
)

func Match(got, want ast.TypeExpr) bool {
	// TODO: right now all traits are empty and so are
	// equivilant to the 'any' type
	if _, ok := want.(*ast.TraitType); ok {
		return true
	}

	switch got := got.(type) {
	case *ast.TraitType:
		// TODO: again all traits are empty for now
		return true
	case *ast.IntType:
		_, ok := want.(*ast.IntType)
		return ok
	case *ast.FloatType:
		_, ok := want.(*ast.FloatType)
		return ok
	case *ast.BoolType:
		_, ok := want.(*ast.BoolType)
		return ok
	case *ast.AtomType:
		_, ok := want.(*ast.AtomType)
		return ok
	case *ast.StringType:
		_, ok := want.(*ast.StringType)
		return ok
	case *ast.PathType:
		_, ok := want.(*ast.PathType)
		return ok
	case *ast.FlagType:
		_, ok := want.(*ast.FlagType)
		return ok
	case *ast.NoneType:
		_, ok := want.(*ast.NoneType)
		return ok
	case *ast.FuncType:
		want, ok := want.(*ast.FuncType)
		if !ok {
			return false
		}

		if len(got.Params.Params) != len(want.Params.Params) {
			return false
		}

		for i := range got.Params.Params {
			if !Match(
				got.Params.Params[i].Type,
				want.Params.Params[i].Type,
			) {
				return false
			}
		}

		if !Match(got.Return, want.Return) {
			return false
		}

		return true
	}

	return false
}

func MatchArity(args []ast.TypeExpr, funcType *ast.FuncType) bool {
	if len(funcType.Params.Params) == 0 && len(args) == 0 {
		return true
	}

	params := funcType.Params.Params
	arity := len(params)
	finalType := params[arity-1].Type
	_, isVariadic := finalType.(*ast.VariadicType)

	if !isVariadic && len(args) == len(funcType.Params.Params) {
		return true
	}

	// params - 1 because the variadic argument can be omitted
	if isVariadic && len(args) > len(funcType.Params.Params)-1 {
		return true
	}

	return false
}

func MatchFunc(args []ast.TypeExpr, funcType *ast.FuncType) bool {
	if !MatchArity(args, funcType) {
		return false
	}

	if len(args) == 0 {
		return true
	}

	params := funcType.Params.Params
	arity := len(params)
	finalType := params[arity-1].Type
	varidicType, ok := finalType.(*ast.VariadicType)
	if ok {
		for i := range args {
			check := varidicType.Type
			if i < len(params)-1 {
				check = params[i].Type
			}
			if !Match(args[i], check) {
				return false
			}
		}

		return true
	}

	for i := range args {
		if !Match(args[i], params[i].Type) {
			return false
		}
	}

	return true
}
