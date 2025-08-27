package thing

import (
	"fmt"
	"log"

	"gortle/internal/ast"
)

type Tag int

const (
	TagTArray    Tag = iota
	TagTList     Tag
	TagTProc     Tag
	TagTPropList Tag
	TagTString   Tag
	TagTSymbol   Tag
	TagTNumber   Tag
	TagTParam    Tag
)

type TArray struct {
	values []*Thing
	size   uint
	dims   uint
	origin uint
}

type TProc struct {
	env    map[TSymbol]*Thing
	tree   ast.ASTNode
	defn   string
	params []TParam
}

type TParam struct {
	name   TSymbol
	dflval *Thing
	opt    bool
	rem    bool
}

type TList []*Thing
type TPropList map[*Thing]*Thing

type TString string
type TNumber float64
type TSymbol string

type Thing struct {
	value interface{}
	tag   Tag
	local bool
}

func New(value interface{}, tag Tag, local bool) *Thing {
	tng := &Thing{
		value: value,
		tag:   tag,
		local: local,
	}
	return tng
}

func NewArray(values []*Thing, size, dims, origin uint, local bool) *Thing {
	arr := TArray{
		values: values,
		size:   size,
		dims:   dims,
		origin: origin,
	}
	return New(arr, TagTArray, local)
}

func NewList(values TList, local bool) *Thing {
	return New(value, TagTList, local)
}

func NewPropList(values TPropList, local bool) *Thing {
	return New(values, TagTPropList, local)
}

func NewProc(env map[TSymbol]*Thing, tree ast.Ast, defn string, params []TParam) *Thing {
	proc := TProc{
		env:    env,
		tree:   tree,
		defn:   defn,
		params: params,
	}
	return New(proc, TagTProc, false)
}

func NewPram(name TSymbol, dflval *Thing, opt, rem bool) *Thing {
	param := TParam{
		name:   name,
		dflval: dflval,
		opt:    opt,
		rem:    rem,
	}
	return New(param, TagTParam, true)
}
