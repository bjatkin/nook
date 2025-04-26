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
	idents map[string]any
}

func (s *scope) lookupIdent(ident string) (any, bool) {
	if value, ok := s.idents[ident]; ok {
		return value, true
	}
	if s.parent != nil {
		return s.parent.lookupIdent(ident)
	}
	return nil, false
}

func (s *scope) setIdent(ident string, value any) {
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
			idents: make(map[string]any),
		},
	}
}

func (vm *VM) WorkingDir() string {
	return vm.workingDir
}

// TODO: send back editor events?
func (vm *VM) Eval(expr ast.Expr) (any, error) {
	switch expr := expr.(type) {
	case ast.SExpr:
		return vm.evalSexpr(expr.Operator, expr.Operands)
	case ast.Int:
		return expr.Value, nil
	case ast.Float:
		return expr.Value, nil
	case ast.Bool:
		return expr.Value, nil
	case ast.String:
		return expr.Value, nil
	case ast.Atom:
		return expr.Value, nil
	case ast.Flag:
		return expr.Value, nil
	case ast.Path:
		return expr.Value, nil
	case ast.Identifier:
		val, ok := vm.scope.lookupIdent(expr.Value.Value)
		if !ok {
			return nil, fmt.Errorf("unknown identifier '%s'", expr.Value.Value)
		}
		return val, nil
	default:
		return nil, fmt.Errorf("invalid expression %#v", expr)
	}
}

func (vm *VM) evalSexpr(operator token.Token, args []ast.Expr) (any, error) {
	switch operator.Kind {
	case token.Plus:
		if len(args) == 0 {
			return nil, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return nil, err
		}

		switch arg := arg.(type) {
		case int64:
			return vm.sumInt(arg, args[1:])
		case float64:
			return vm.sumFloat(arg, args[1:])
		default:
			return nil, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Minus:
		if len(args) == 0 {
			return nil, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return nil, err
		}

		switch arg := arg.(type) {
		case int64:
			return vm.minusInt(arg, args[1:])
		case float64:
			return vm.minusFloat(arg, args[1:])
		default:
			return nil, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Divide:
		if len(args) == 0 {
			return nil, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return nil, err
		}

		switch arg := arg.(type) {
		case int64:
			return vm.divInt(arg, args[1:])
		case float64:
			return vm.divFloat(arg, args[1:])
		default:
			return nil, fmt.Errorf("'%s' must have type of float or int", operator.Value)
		}
	case token.Multiply:
		if len(args) == 0 {
			return nil, nil
		}

		arg, err := vm.Eval(args[0])
		if err != nil {
			return nil, err
		}

		switch arg := arg.(type) {
		case int64:
			return vm.mulInt(arg, args[1:])
		case float64:
			return vm.mulFloat(arg, args[1:])
		default:
			return nil, fmt.Errorf("'%v' must have type of float or int", arg)
		}
	case token.Let:
		if len(args) != 2 {
			return nil, fmt.Errorf("must have exactly 2 arguments for let")
		}
		ident, ok := args[0].(ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("first argument must be an identifier")
		}
		value, err := vm.Eval(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid value for let assignment")
		}

		vm.scope.setIdent(ident.Value.Value, value)
		return nil, nil
	case token.Identifier:
		switch operator.Value {
		case "cd":
			if len(args) != 1 {
				return nil, fmt.Errorf("cd takes exactly 1 argument")
			}

			dir, err := vm.Eval(args[0])
			if err != nil {
				return nil, fmt.Errorf("failed to evali first argument")
			}

			strDir, ok := dir.(string)
			if !ok {
				return nil, fmt.Errorf("'%v' is not a valid path", dir)
			}

			workingDir := strDir
			if !strings.HasPrefix(strDir, "/") {
				workingDir = path.Join(vm.workingDir, strDir)
			}

			// make sure the directory exists before switching to it
			if _, err := os.Stat(workingDir); err != nil {
				return nil, fmt.Errorf("'%v' was not found: '%w'", workingDir, err)
			}

			vm.workingDir = workingDir
			return nil, nil
		case "ls":
			// TODO: add expected arguments like -la
			if len(args) > 0 {
				return nil, fmt.Errorf("'ls' takes zero arguments")
			}

			files, err := os.ReadDir(vm.workingDir)
			if err != nil {
				return nil, fmt.Errorf("could not read dir '%v': '%w'", vm.workingDir, err)
			}

			// TODO: hand back a list of datastructures here so that nook-script can interact
			// with the returned value
			found := []string{}
			for _, file := range files {
				found = append(found, file.Name())
			}

			return found, nil
		case "ex":
			if len(args) == 0 {
				return nil, fmt.Errorf("'ex' requires a program to run")
			}

			nameValue, err := vm.Eval(args[0])
			if err != nil {
				return nil, fmt.Errorf("failed to get command name: '%w'", err)
			}

			// TODO: support atoms here as well
			name, ok := nameValue.(string)
			if !ok {
				return nil, fmt.Errorf("command name must be a string but got '%v'", nameValue)
			}
			// TODO: this will break for strings in the form ":test" so we need to fix that. Once we
			// introcue a proper value systme for the vm that will help fix this
			name = strings.TrimPrefix(name, ":")

			cmdArgs := []string{}
			// TODO: this logic can probably be centralized somewhere but we'll just do it here for now
			for _, arg := range args[1:] {
				cmdArg, err := vm.Eval(arg)
				if err != nil {
					return nil, fmt.Errorf("failed to execute argument '%v'", arg)
				}

				// TODO: support atoms, paths, strings, and flags here
				// That won't matter until we get a proper value system in the vm though
				// because for now all those types are just represented by go strings
				strArg, ok := cmdArg.(string)
				if !ok {
					return nil, fmt.Errorf("argument must be a string: '%v'", cmdArg)
				}
				strArg = strings.TrimPrefix(strArg, ":")

				cmdArgs = append(cmdArgs, strArg)
			}

			cmd := exec.Command(name, cmdArgs...)
			cmd.Dir = vm.workingDir

			// we can ignore the error here since a failed call isn't a vm failue and should
			// be represented by the value.
			// TODO: need to set the status code for the last command
			result, _ := cmd.CombinedOutput()
			return string(result), nil
		default:
			return nil, fmt.Errorf("invalid operator: '%s'", operator.Value)
		}
	default:
		return nil, fmt.Errorf("'%v' is not a valid s-expr operator", operator.Value)
	}
}

func (vm *VM) sumInt(start int64, args []ast.Expr) (int64, error) {
	sum := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		i, ok := value.(int64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be an integer", value)
		}
		sum += i
	}

	return sum, nil
}

func (vm *VM) sumFloat(start float64, args []ast.Expr) (float64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		f, ok := value.(float64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be a float", value)
		}
		ret += f
	}

	return ret, nil
}

func (vm *VM) minusInt(start int64, args []ast.Expr) (int64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		i, ok := value.(int64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be an int", value)
		}
		ret -= i
	}

	return ret, nil
}

func (vm *VM) minusFloat(start float64, args []ast.Expr) (float64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		f, ok := value.(float64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be a float", value)
		}
		ret -= f
	}

	return ret, nil
}

func (vm *VM) divInt(start int64, args []ast.Expr) (int64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		i, ok := value.(int64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be an int", value)
		}
		ret /= i
	}

	return ret, nil
}

func (vm *VM) divFloat(start float64, args []ast.Expr) (float64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		f, ok := value.(float64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be a float", value)
		}
		ret /= f
	}

	return ret, nil
}

func (vm *VM) mulInt(start int64, args []ast.Expr) (int64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		i, ok := value.(int64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be an int", value)
		}
		ret *= i
	}

	return ret, nil
}

func (vm *VM) mulFloat(start float64, args []ast.Expr) (float64, error) {
	ret := start
	for _, arg := range args {
		value, err := vm.Eval(arg)
		if err != nil {
			return 0, err
		}

		f, ok := value.(float64)
		if !ok {
			return 0, fmt.Errorf("'%v' must be a float", value)
		}
		ret *= f
	}

	return ret, nil
}
