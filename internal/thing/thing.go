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

	defaultListSize uint = 32
)

type TArray struct {
	values []*Thing
	size   uint
	dims   uint
	origin uint
}

type TProc struct {
	env    TEnv
	tree   ast.ASTNode
	defn   string
	params []*Thing
}

type TParam struct {
	name   TSymbol
	dflval *Thing
	opt    bool
	rem    bool
}

type TList []*Thing
type TPropList map[*Thing]*Thing
type TEnv map[TSymbol]*Thing

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

func NewArray(size, dims, origin uint, local bool) *Thing {
	arr := TArray{
		values: make([]*Thing, 0, size),
		size:   size,
		dims:   dims,
		origin: origin,
	}
	return New(arr, TagTArray, local)
}

func NewList(local bool) *Thing {
	return New(make(TList, 0, defaultListSize), TagTList, local)
}

func NewPropList(local bool) *Thing {
	return New(make(TList, 0, 1024), TagTPropList, local)
}

func NewProc(tree ast.Ast, defn string, params []TParam) *Thing {
	proc := TProc{
		env:    make(TEnv, defaultListSize),
		tree:   tree,
		defn:   defn,
		params: params,
	}
	return New(proc, TagTProc, false)
}

func NewParam(name TSymbol, dflval *Thing, opt, rem bool) *Thing {
	param := TParam{
		name:   name,
		dflval: dflval,
		opt:    opt,
		rem:    rem,
	}
	return New(param, TagTParam, true)
}
