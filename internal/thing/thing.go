package thing

import (
	"fmt"
	"log"

	"gortle/internal/ast"
)

type Tag int

const (
	TagTArray  Tag = iota
	TagTList   Tag
	TagTProc   Tag
	TagTEnv    Tag
	TagTString Tag
	TagTSymbol Tag
	TagTNumber Tag
)

type TArray struct {
	values []*Thing
	size   uint
	dims   uint
	origin uint
}

type TList struct {
	values []*Thing
	size   uint
}

type TProc struct {
	env     TEnv
	tree    ast.ASTNode
	nparams int
}

type TEnv map[TSymbol]*Thing

type TString string
type TNumber float64
type TSymbol string

type Thing struct {
	value interface{}
	tag   Tag
}
