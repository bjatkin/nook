package checker

import (
	"fmt"
	"strings"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/builtin"
	"github.com/bjatkin/nook/script/symbol"
	"github.com/bjatkin/nook/script/types"
)

type Checker struct {
	table  *symbol.Table
	Errors []error
}

func NewChecker() *Checker {
	table := symbol.NewTable()

	// add builtins to the symbole table
	for _, builtin := range builtin.Builtins {
		table.AddBuiltin(builtin)
	}

	return &Checker{
		table: table,
	}
}

func (c *Checker) openScope() {
	scope := c.table.OpenScope()
	c.table = scope

}

func (c *Checker) closeScope() {
	// TODO: check for nil parent
	c.table = c.table.CloseScope()
}

func (c *Checker) addError(err error) {
	c.Errors = append(c.Errors, err)
}

// Infer infers types for all expressions to prepare for type checking
func (c *Checker) Infer(expr ast.Expr) ast.TypeExpr {
	switch expr := expr.(type) {
	case *ast.Int:
		return &ast.IntType{}
	case *ast.Float:
		return &ast.FloatType{}
	case *ast.BoolType:
		return &ast.BoolType{}
	case *ast.Atom:
		return &ast.AtomType{}
	case *ast.String:
		return &ast.StringType{}
	case *ast.Path:
		return &ast.PathType{}
	case *ast.Flag:
		return &ast.FlagType{}
	case *ast.Bool:
		return &ast.BoolType{}
	case *ast.Command:
		return &ast.CommandType{}
	case *ast.Func:
		c.openScope()
		bodyType := c.Infer(expr.Body)
		if !types.Match(bodyType, expr.Type.Return) {
			c.addError(fmt.Errorf("body type does not expected function return type '%v'", bodyType))
		}
		c.closeScope()

		return expr.Type
	case *ast.Let:
		// TODO: handle scope
		c.openScope()
		exprType := c.Infer(expr.Value)
		c.closeScope()

		err := c.table.AddLet(expr, exprType)
		if err != nil {
			c.addError(err)
		}

		// let expressions return a none value
		return &ast.NoneType{}
	case *ast.Identifier:
		identEntry, ok := c.table.LookupValue(expr.Name)
		if !ok {
			c.addError(fmt.Errorf("identifier '%s' has not been defined", expr.Name))
			return &ast.NoneType{}
		}

		return identEntry.Type
	case *ast.Call:
		return c.inferCall(expr)
	default:
		panic(fmt.Sprintf("failed to infer type '%v'", expr))
	}
}

func (c *Checker) inferCall(call *ast.Call) ast.TypeExpr {
	c.openScope()
	defer c.closeScope()

	argTypes := []ast.TypeExpr{}
	for _, arg := range call.Args {
		argType := c.Infer(arg)
		argTypes = append(argTypes, argType)
	}

	switch expr := call.Func.(type) {
	case *ast.Func:
		typeExpr := c.Infer(expr)

		funcType := typeExpr.(*ast.FuncType)

		if len(argTypes) != len(funcType.Params.Params) {
			c.addError(fmt.Errorf("arities do not match"))
		}

		return expr.Type.Return
	case *ast.Identifier:
		entry, ok := c.table.Lookup(expr.Name)
		if !ok {
			c.addError(fmt.Errorf("unknown identifier '%v'", expr.Name))
		}

		switch entry := entry.(type) {
		case *symbol.ValueEntry:
			return c.checkFuncCall(entry.Type, argTypes)
		case *symbol.ImplEntry:
			impl, ok := c.checkImplCall(entry, argTypes)
			if !ok {
				return &ast.NoneType{}
			}

			// swap the operator out for the actual function
			call.Func = impl.Decl.Func

			return impl.Type.Return
		case *symbol.BuiltinEntry:
			builtin, ok := c.checkBuiltinCall(entry, argTypes)
			if !ok {
				return &ast.NoneType{}
			}

			// swap the operator out for the correct builtin
			call.Func = builtin.Decl

			return builtin.Type.Return
		default:
			panic("invalid symbole table entry")
		}

	default:
		c.addError(fmt.Errorf("can not call value with type '%v'", expr))
		return &ast.NoneType{}
	}
}

func (c *Checker) checkBuiltinCall(builtin *symbol.BuiltinEntry, args []ast.TypeExpr) (*symbol.BuiltinOverload, bool) {
	overload, ok := builtin.Match(args)
	if !ok {
		argTypes := []string{}
		for i := range args {
			switch args[i].(type) {
			case *ast.IntType:
				argTypes = append(argTypes, "int")
			case *ast.FloatType:
				argTypes = append(argTypes, "float")
			case *ast.BoolType:
				argTypes = append(argTypes, "bool")
			case *ast.FlagType:
				argTypes = append(argTypes, "flag")
			case *ast.CommandType:
				argTypes = append(argTypes, "command")
			case *ast.PathType:
				argTypes = append(argTypes, "path")
			case *ast.FuncType:
				argTypes = append(argTypes, "func")
			case *ast.TraitType:
				// TODO: this should be more specific than just an any
				argTypes = append(argTypes, "any")
			}
		}
		c.addError(fmt.Errorf("could not find a matching overload for ('%s' %s)", builtin.Name, strings.Join(argTypes, " ")))
		return nil, false
	}

	return overload, true
}

func (c *Checker) checkImplCall(impl *symbol.ImplEntry, args []ast.TypeExpr) (*symbol.ImplOverload, bool) {
	overload, ok := impl.Match(args)
	if !ok {
		c.addError(fmt.Errorf("could not find a matching overload"))
		return nil, false
	}

	return overload, true
}

func (c *Checker) checkFuncCall(typeExpr ast.TypeExpr, args []ast.TypeExpr) ast.TypeExpr {
	funcType, ok := typeExpr.(*ast.FuncType)
	if !ok {
		c.addError(fmt.Errorf("call type is not a function type %v", funcType))
	}

	if funcType == nil {
		c.addError(fmt.Errorf("invalid function call"))
		return nil
	}

	if funcType.Params == nil && len(args) != 0 {
		c.addError(fmt.Errorf("expected no arguments but got %d", len(args)))
		return funcType.Return
	}

	if len(args) != len(funcType.Params.Params) {
		c.addError(fmt.Errorf("arities do not match"))
		return funcType.Return
	}

	for i, arg := range args {
		wantType := funcType.Params.Params[i].Type
		if !types.Match(arg, wantType) {
			c.addError(fmt.Errorf("argument type is incorrect got '%v' but wanted '%v'", arg, wantType))
		}
	}

	return funcType.Return
}
