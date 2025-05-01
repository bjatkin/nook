package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

type scope struct {
	parent *scope
	idents map[string]Value
}

func (s *scope) lookupIdent(ident string) (Value, bool) {
	if value, ok := s.idents[ident]; ok {
		return value, true
	}
	if s.parent != nil {
		return s.parent.lookupIdent(ident)
	}
	return Value{}, false
}

func (s *scope) setIdent(ident string, value Value) {
	s.idents[ident] = value
}

type VM struct {
	workingDir string
	scope      *scope
}

func NewVM(workingDir string) VM {
	return VM{
		workingDir: workingDir,
		scope: &scope{
			idents: make(map[string]Value),
		},
	}
}

func (vm *VM) WorkingDir() string {
	return vm.workingDir
}

// TODO: send back editor events?
func (vm *VM) Eval(expr ast.Expr) (Value, error) {
	switch expr := expr.(type) {
	case ast.SExpr:
		return vm.evalSexpr(expr.Operator, expr.Operands)
	case ast.Int:
		return Value{value: expr.Value, kind: Int}, nil
	case ast.Float:
		return Value{value: expr.Value, kind: Float}, nil
	case ast.Bool:
		return Value{value: expr.Value, kind: Bool}, nil
	case ast.String:
		return Value{value: expr.Value, kind: String}, nil
	case ast.Atom:
		return Value{value: expr.Value, kind: Atom}, nil
	case ast.Flag:
		return Value{value: expr.Value, kind: Flag}, nil
	case ast.Path:
		return Value{value: expr.Value, kind: Path}, nil
	case ast.Identifier:
		val, ok := vm.scope.lookupIdent(expr.Value.Value)
		if !ok {
			return Value{}, fmt.Errorf("unknown identifier '%s'", expr.Value.Value)
		}
		return val, nil
	default:
		return Value{}, fmt.Errorf("invalid expression %#v", expr)
	}
}

func (vm *VM) evalSexpr(operator token.Token, args []ast.Expr) (Value, error) {
	switch operator.Kind {
	case token.Plus:
		if len(args) == 0 {
			return Value{}, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return Value{}, err
		}

		switch arg.kind {
		case Int:
			return vm.sumInt(arg.Int(), args[1:])
		case Float:
			return vm.sumFloat(arg.Float(), args[1:])
		default:
			return Value{}, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Minus:
		if len(args) == 0 {
			return Value{}, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return Value{}, err
		}

		switch arg.kind {
		case Int:
			return vm.minusInt(arg.Int(), args[1:])
		case Float:
			return vm.minusFloat(arg.Float(), args[1:])
		default:
			return Value{}, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Divide:
		if len(args) == 0 {
			return Value{}, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return Value{}, err
		}

		switch arg.kind {
		case Int:
			return vm.divInt(arg.Int(), args[1:])
		case Float:
			return vm.divFloat(arg.Float(), args[1:])
		default:
			return Value{}, fmt.Errorf("'%s' must have type of float or int", operator.Value)
		}
	case token.Multiply:
		if len(args) == 0 {
			return Value{}, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return Value{}, err
		}

		switch arg.kind {
		case Int:
			return vm.mulInt(arg.Int(), args[1:])
		case Float:
			return vm.mulFloat(arg.Float(), args[1:])
		default:
			return Value{}, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Let:
		if len(args) != 2 {
			return Value{}, fmt.Errorf("must have exactly 2 arguments for let")
		}
		ident, ok := args[0].(ast.Identifier)
		if !ok {
			return Value{}, fmt.Errorf("first argument must be an identifier")
		}
		value, err := vm.Eval(args[1])
		if err != nil {
			return Value{}, fmt.Errorf("invalid value for let assignment")
		}

		vm.scope.setIdent(ident.Value.Value, value)
		return Value{}, nil
	case token.Command:
		name := operator.Value[1:]
		cmdArgs := []string{}
		for _, arg := range args {
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
		cmd.Dir = vm.workingDir
		result, _ := cmd.CombinedOutput()

		// TODO: handle status codes and return a more complex type than this
		return Value{value: string(result), kind: String}, nil
	case token.Identifier:
		switch operator.Value {
		case "cd":
			if len(args) != 1 {
				return Value{}, fmt.Errorf("cd takes exactly 1 argument")
			}

			dir, err := vm.Eval(args[0])
			if err != nil {
				return Value{}, fmt.Errorf("failed to evali first argument")
			}

			if dir.kind != String {
				return Value{}, fmt.Errorf("'%v' is not a valid path", dir)
			}

			strDir := dir.String()
			workingDir := strDir
			if !strings.HasPrefix(strDir, "/") {
				workingDir = path.Join(vm.workingDir, strDir)
			}

			// make sure the directory exists before switching to it
			if _, err := os.Stat(workingDir); err != nil {
				return Value{}, fmt.Errorf("'%v' was not found: '%w'", workingDir, err)
			}

			vm.workingDir = workingDir
			return NoneValue, nil
		case "ls":
			// TODO: add expected arguments like -la?
			if len(args) > 0 {
				return Value{}, fmt.Errorf("'ls' takes zero arguments")
			}

			files, err := os.ReadDir(vm.workingDir)
			if err != nil {
				return Value{}, fmt.Errorf("could not read dir '%v': '%w'", vm.workingDir, err)
			}

			// TODO: hand back a list of datastructures here so that nook-script can interact
			// with the returned value
			found := []string{}
			for _, file := range files {
				found = append(found, file.Name())
			}

			return Value{value: found, kind: String}, nil
		default:
			return Value{}, fmt.Errorf("invalid operator: '%s'", operator.Value)
		}
	default:
		return Value{}, fmt.Errorf("'%v' is not a valid s-expr operator", operator.Value)
	}
}

func (vm *VM) evalArgs(kind Kind, args []ast.Expr) ([]Value, error) {
	evaled := []Value{}
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return nil, err
		}

		if value.kind != kind {
			return nil, fmt.Errorf("%s(%v) must be an %s", value.kind, value.value, kind)
		}
		evaled = append(evaled, value)
	}

	return evaled, nil
}

func (vm *VM) sumInt(start int64, args []ast.Expr) (Value, error) {
	ints, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	sum := start
	for _, i := range ints {
		sum += i.Int()
	}

	return Value{value: sum, kind: Int}, nil
}

func (vm *VM) minusInt(start int64, args []ast.Expr) (Value, error) {
	ints, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, i := range ints {
		min -= i.Int()
	}

	return Value{value: min, kind: Int}, nil
}

func (vm *VM) divInt(start int64, args []ast.Expr) (Value, error) {
	ints, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, i := range ints {
		min /= i.Int()
	}

	return Value{value: min, kind: Int}, nil
}

func (vm *VM) mulInt(start int64, args []ast.Expr) (Value, error) {
	ints, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, i := range ints {
		min *= i.Int()
	}

	return Value{value: min, kind: Int}, nil
}

func (vm *VM) sumFloat(start float64, args []ast.Expr) (Value, error) {
	floats, err := vm.evalArgs(Float, args)
	if err != nil {
		return Value{}, err
	}

	sum := start
	for _, f := range floats {
		sum += f.Float()
	}

	return Value{value: sum, kind: Float}, nil
}

func (vm *VM) minusFloat(start float64, args []ast.Expr) (Value, error) {
	floats, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, f := range floats {
		min -= f.Float()
	}

	return Value{value: min, kind: Int}, nil
}

func (vm *VM) divFloat(start float64, args []ast.Expr) (Value, error) {
	floats, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, f := range floats {
		min /= f.Float()
	}

	return Value{value: min, kind: Int}, nil
}

func (vm *VM) mulFloat(start float64, args []ast.Expr) (Value, error) {
	floats, err := vm.evalArgs(Int, args)
	if err != nil {
		return Value{}, err
	}

	min := start
	for _, f := range floats {
		min *= f.Float()
	}

	return Value{value: min, kind: Int}, nil
}
