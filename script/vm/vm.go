package vm

import (
	"fmt"
	"os/exec"

	"github.com/bjatkin/nook/script/ast"
)

type scope struct {
	parent *scope
	idents map[string]Value
}

func (s *scope) lookupIdent(ident string, args []ast.Expr) (Value, bool) {
	if value, ok := s.idents[ident]; ok {
		return value, true
	}

	if s.parent == nil {
		return Value{}, false
	}

	return s.parent.lookupIdent(ident, args)
}

func (s *scope) setIdent(ident string, value Value) {
	s.idents[ident] = value
}

type VM struct {
	scope *scope
}

func NewVM() *VM {
	return &VM{
		scope: &scope{
			idents: make(map[string]Value),
		},
	}
}

// TODO: send back editor events?
func (vm *VM) Eval(expr ast.Expr) (Value, error) {
	switch expr := expr.(type) {
	case *ast.Let:
		value, err := vm.Eval(expr.Value)
		if err != nil {
			return Value{}, fmt.Errorf("failed to eval let expr")
		}
		vm.scope.setIdent(expr.Identifier.Name, value)

		return NoneValue, nil
	case *ast.Call:
		value, err := vm.evalCall(expr.Func, expr.Args)
		if err != nil {
			return Value{}, fmt.Errorf("failed to call expr: '%w'", err)
		}

		return value, nil
	case *ast.Command:
		name := expr.Name
		cmdArgs := []string{}
		for _, arg := range expr.Args {
			value, err := vm.Eval(arg)
			if err != nil {
				return Value{}, fmt.Errorf("failed to eval argument: %v", arg)
			}

			// TODO: really need an actual value type, not just any
			// also, traits are how we should do this
			// support anything that can be turnned into a shell value
			// BUT, for now just turn everything into a string
			strValue := value.String()
			if value.kind == Atom {
				strValue = strValue[1:]
			}

			cmdArgs = append(cmdArgs, strValue)
		}

		cmd := exec.Command(name, cmdArgs...)
		result, _ := cmd.CombinedOutput()

		// TODO: handle status codes and return a more complex type than this
		return Value{value: string(result), kind: String}, nil
	case *ast.Int:
		return Value{value: expr.Value, kind: Int}, nil
	case *ast.Float:
		return Value{value: expr.Value, kind: Float}, nil
	case *ast.Bool:
		return Value{value: expr.Value, kind: Bool}, nil
	case *ast.String:
		return Value{value: expr.Value, kind: String}, nil
	case *ast.Atom:
		return Value{value: expr.Value, kind: Atom}, nil
	case *ast.Flag:
		return Value{value: expr.Value, kind: Flag}, nil
	case *ast.Path:
		return Value{value: expr.Value, kind: Path}, nil
	case *ast.Nil:
		return Value{value: nil, kind: None}, nil
	case *ast.Identifier:
		val, ok := vm.scope.lookupIdent(expr.Name, nil)
		if !ok {
			return Value{}, fmt.Errorf("unknown identifier '%s'", expr.Name)
		}
		return val, nil
	default:
		return Value{}, fmt.Errorf("invalid runtime expression %#v", expr)
	}
}

func (vm *VM) evalCall(operator ast.Expr, args []ast.Expr) (Value, error) {
	switch operator := operator.(type) {
	case *ast.Func:
		panic("ast.Func is not currently supported")
	case *ast.Identifier:
		panic("ast.Identifier is not currently supported")
	case *ast.Builtin:
		values, err := vm.evalArgs(args)
		if err != nil {
			return Value{}, fmt.Errorf("failed to evaluate argument '%w'", err)
		}

		args := []any{}
		for i := range values {
			args = append(args, values[i].value)
		}

		ret, err := operator.Fn(args...)
		if err != nil {
			return Value{}, err
		}
		if ret == nil {
			return NoneValue, nil
		}

		switch ret.(type) {
		case int64:
			return Value{value: ret, kind: Int}, nil
		case float64:
			return Value{value: ret, kind: Float}, nil
		case string:
			return Value{value: ret, kind: String}, nil
		default:
			return Value{}, fmt.Errorf("failed to convert return to return type")
		}

	default:
		panic("invalid call identifier")
	}
}

func (vm *VM) evalArgs(args []ast.Expr) ([]Value, error) {
	evaled := []Value{}
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return nil, err
		}

		evaled = append(evaled, value)
	}

	return evaled, nil
}
