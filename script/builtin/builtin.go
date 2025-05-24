package builtin

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bjatkin/nook/script/ast"
)

// Builtin is a builtin function for the nook environment
type Builtin struct {
	Name string
	Type *ast.FuncType
	Fn   func(args ...any) (any, error)
}

// Builtins is a slice of all the nook builtin functions.
var Builtins = []Builtin{
	{
		Name: "+",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.IntType{}}}},
			},
			Return: &ast.IntType{},
		},
		Fn: AddInt64,
	},
	{
		Name: "+",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.FloatType{}}}},
			},
			Return: &ast.FloatType{},
		},
		Fn: AddFloat64,
	},
	{
		Name: "-",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.IntType{}}}},
			},
			Return: &ast.IntType{},
		},
		Fn: MinusInt64,
	},
	{
		Name: "-",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.FloatType{}}}},
			},
			Return: &ast.FloatType{},
		},
		Fn: MinusFloat64,
	},
	{
		Name: "*",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.IntType{}}}},
			},
			Return: &ast.IntType{},
		},
		Fn: MultiplyInt64,
	},
	{
		Name: "*",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.FloatType{}}}},
			},
			Return: &ast.FloatType{},
		},
		Fn: MultiplyFloat64,
	},
	{
		Name: "/",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.IntType{}}}},
			},
			Return: &ast.IntType{},
		},
		Fn: DivideInt64,
	},
	{
		Name: "/",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.VariadicType{Type: &ast.FloatType{}}}},
			},
			Return: &ast.FloatType{},
		},
		Fn: DivideFloat64,
	},
	{
		Name: "cd",
		Type: &ast.FuncType{
			Params: &ast.ParamList{
				Params: []ast.Param{{Type: &ast.PathType{}}},
			},
			// TODO: this should really return an option type that's either error or none, but I'm not sure how to do that yet...
			Return: &ast.NoneType{},
		},
		Fn: ChangeDir,
	},
	{
		Name: "ls",
		Type: &ast.FuncType{
			Params: &ast.ParamList{},
			// TODO: sould be a slice of data structures
			Return: &ast.StringType{},
		},
		Fn: ListFiles,
	},
}

// builtin functions
var AddInt64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	sum := args[0].(int64)
	for _, arg := range args[1:] {
		sum += arg.(int64)
	}
	return sum, nil
}

var AddFloat64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return 0, nil
	}

	sum := args[0].(float64)
	for _, arg := range args[1:] {
		sum += arg.(float64)
	}
	return sum, nil
}

var MinusInt64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	min := args[0].(int64)
	for _, arg := range args[1:] {
		min -= arg.(int64)
	}
	return min, nil
}

var MinusFloat64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return float64(0), nil
	}

	min := args[0].(float64)
	for _, arg := range args[1:] {
		min -= arg.(float64)
	}
	return min, nil
}

var MultiplyInt64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	min := args[0].(int64)
	for _, arg := range args[1:] {
		min *= arg.(int64)
	}
	return min, nil
}

var MultiplyFloat64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return float64(0), nil
	}

	min := args[0].(float64)
	for _, arg := range args[1:] {
		min *= arg.(float64)
	}
	return min, nil
}

var DivideInt64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	min := args[0].(int64)
	for _, arg := range args[1:] {
		min /= arg.(int64)
	}
	return min, nil
}

var DivideFloat64 = func(args ...any) (any, error) {
	if len(args) == 0 {
		return float64(0), nil
	}

	min := args[0].(float64)
	for _, arg := range args[1:] {
		min /= arg.(float64)
	}
	return min, nil
}

var ChangeDir = func(args ...any) (any, error) {
	dir := args[0].(string)

	// make sure the directory exists before switching to it
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("'%v' was not found: '%v'", dir, err)
	}

	if !strings.HasPrefix(dir, "/") {
		workingDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory")
		}

		dir = path.Join(workingDir, dir)
	}

	err := os.Chdir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to change directories '%w'", err)
	}

	return nil, nil
}

var ListFiles = func(args ...any) (any, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory '%v'", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read dir '%s': '%v'", dir, err)
	}

	// TODO: hand back a list of data structures here so that nook-script can interact
	// with the returned value rather than just getting a list of strings
	found := []string{}
	for _, file := range files {
		found = append(found, "\""+file.Name()+"\"")
	}

	return "[ " + strings.Join(found, " ") + " ]", nil
}
