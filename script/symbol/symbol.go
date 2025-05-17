package symbol

import (
	"fmt"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/builtin"
	"github.com/bjatkin/nook/script/types"
)

type entryKind int

const (
	valueKind = entryKind(0)
	implKind
	builtinKind
)

type Entry interface {
	kind() entryKind
}

type ValueEntry struct {
	Name string
	Type ast.TypeExpr
	Decl *ast.Let
}

func (e *ValueEntry) kind() entryKind { return valueKind }

type ImplOverload struct {
	Type *ast.FuncType
	Decl *ast.Impl
}

type ImplEntry struct {
	Name      string
	Overloads []ImplOverload
}

func (e *ImplEntry) kind() entryKind { return implKind }

func (e *ImplEntry) Match(args []ast.TypeExpr) (*ImplOverload, bool) {
	for _, overload := range e.Overloads {
		if types.MatchFunc(args, overload.Type) {
			return &overload, true
		}
	}

	return nil, false
}

type BuiltinOverload struct {
	Type *ast.FuncType
	Decl *ast.Builtin
}

type BuiltinEntry struct {
	Name      string
	Overloads []BuiltinOverload
}

func (e *BuiltinEntry) kind() entryKind { return builtinKind }

func (e *BuiltinEntry) Match(args []ast.TypeExpr) (*BuiltinOverload, bool) {
	for _, overload := range e.Overloads {
		if types.MatchFunc(args, overload.Type) {
			return &overload, true
		}
	}

	return nil, false
}

type Table struct {
	parent   *Table
	symboles map[string]Entry
}

func NewTable() *Table {
	return &Table{
		symboles: map[string]Entry{},
	}
}

func (t *Table) AddBuiltin(builtin builtin.Builtin) error {
	name := builtin.Name
	entry, ok := t.symboles[name]
	if !ok {
		t.addNewBuiltin(builtin)
		return nil
	}

	builtinEntry, ok := entry.(*BuiltinEntry)
	if !ok {
		return fmt.Errorf("entry already exists and is not a builtin '%v'", entry)
	}

	builtinEntry.Overloads = append(builtinEntry.Overloads, BuiltinOverload{
		Type: builtin.Type,
		Decl: &ast.Builtin{
			Fn: builtin.Fn,
		},
	})

	t.symboles[name] = builtinEntry
	return nil
}

func (t *Table) addNewBuiltin(builtin builtin.Builtin) {
	name := builtin.Name
	entry := &BuiltinEntry{
		Name: name,
		Overloads: []BuiltinOverload{{
			Type: builtin.Type,
			Decl: &ast.Builtin{
				Fn: builtin.Fn,
			},
		}},
	}

	t.symboles[name] = entry
}

func (t *Table) AddImpl(impl *ast.Impl) error {
	name := impl.Identifier.Name
	entry, ok := t.symboles[name]
	if !ok {
		t.addNewImpl(impl)
		return nil
	}

	implEntry, ok := entry.(*ImplEntry)
	if !ok {
		return fmt.Errorf("entry already exists and is not an impl '%v'", entry)
	}

	implEntry.Overloads = append(implEntry.Overloads, ImplOverload{
		Type: impl.Func.Type,
		Decl: impl,
	})

	t.symboles[name] = implEntry
	return nil
}

func (t *Table) addNewImpl(impl *ast.Impl) {
	name := impl.Identifier.Name
	entry := &ImplEntry{
		Name: name,
		Overloads: []ImplOverload{{
			Type: impl.Func.Type,
			Decl: impl,
		}},
	}

	t.symboles[name] = entry
}

func (t *Table) AddLet(let *ast.Let, letType ast.TypeExpr) error {
	name := let.Identifier.Name
	entry, ok := t.symboles[name]
	if !ok {
		t.addNewLet(let, letType)
		return nil
	}

	_, ok = entry.(*ValueEntry)
	if !ok {
		return fmt.Errorf("entry alread exists and is not a value '%v'", entry)
	}

	t.symboles[name] = &ValueEntry{
		Name: name,
		Type: letType,
		Decl: let,
	}
	return nil
}

func (t *Table) addNewLet(let *ast.Let, letType ast.TypeExpr) {
	name := let.Identifier.Name
	t.symboles[name] = &ValueEntry{
		Name: name,
		Type: letType,
		Decl: let,
	}
}

func (t *Table) Lookup(name string) (Entry, bool) {
	valueEntry, ok := t.lookupValueInScope(name)
	if ok {
		return valueEntry, ok
	}

	implEntry, ok := t.lookupImplInScope(name)
	if ok {
		return implEntry, ok
	}

	builtinEntry, ok := t.lookupBuiltinInScope(name)
	if ok {
		return builtinEntry, ok
	}

	if t.parent == nil {
		return nil, false
	}

	return t.parent.Lookup(name)
}

func (t *Table) LookupValue(name string) (*ValueEntry, bool) {
	if entry, ok := t.lookupValueInScope(name); ok {
		return entry, true
	}

	if t.parent == nil {
		return nil, false
	}

	return t.parent.LookupValue(name)
}

func (t *Table) lookupValueInScope(name string) (*ValueEntry, bool) {
	entry, ok := t.symboles[name]
	if !ok {
		return nil, false
	}

	valueEntry, ok := entry.(*ValueEntry)
	if !ok {
		return nil, false
	}

	return valueEntry, true
}

func (t *Table) LookupImpl(name string) (*ImplEntry, bool) {
	if entry, ok := t.lookupImplInScope(name); ok {
		return entry, true
	}

	if t.parent == nil {
		return nil, false
	}

	return t.parent.LookupImpl(name)
}

func (t *Table) lookupImplInScope(name string) (*ImplEntry, bool) {
	entry, ok := t.symboles[name]
	if !ok {
		return nil, false
	}

	implEntry, ok := entry.(*ImplEntry)
	if !ok {
		return nil, false
	}

	return implEntry, true
}

func (t *Table) LookupBuiltin(name string) (*BuiltinEntry, bool) {
	if entry, ok := t.lookupBuiltinInScope(name); ok {
		return entry, true
	}

	if t.parent == nil {
		return nil, false
	}

	return t.parent.LookupBuiltin(name)
}

func (t *Table) lookupBuiltinInScope(name string) (*BuiltinEntry, bool) {
	entry, ok := t.symboles[name]
	if !ok {
		return nil, false
	}

	implEntry, ok := entry.(*BuiltinEntry)
	if !ok {
		return nil, false
	}

	return implEntry, true
}

func (t *Table) OpenScope() *Table {
	scope := NewTable()
	scope.parent = t
	return scope
}

func (t *Table) CloseScope() *Table {
	return t.parent
}
