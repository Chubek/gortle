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
	dims   []uint
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

func NewArray(dims []uint, origin uint, local bool) *Thing {
	if len(dims) == 0 {
		log.Fatalf("newarray: dimens are zero")
		return
	}

	size := 1
	for _, dim := range dims {
		size *= dim
	}

	values := make([]*Thing, size)
	arr := &TArray{
		values: values,
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
	proc := &TProc{
		env:    make(TEnv, defaultListSize),
		tree:   tree,
		defn:   defn,
		params: params,
	}
	return New(proc, TagTProc, false)
}

func NewParam(name TSymbol, dflval *Thing, opt, rem bool) *Thing {
	param := &TParam{
		name:   name,
		dflval: dflval,
		opt:    opt,
		rem:    rem,
	}
	return New(param, TagTParam, true)
}

func NewString(value TString, local bool) *Thing {
	return New(value, TagTString, local)
}

func NewNumber(value TNumber, local bool) *Thing {
	return New(value, TagTNumber, local)
}

func NewSymbol(value TSymbol, local bool) *Thing {
	return New(value, TagTSymbol, local)
}

func (t *TArray) getIndex(coords []int) (int, error) {
	if len(coords) != len(t.dims) {
		return 0, fmt.Errorf("Invalid number of coordinates")
	}

	index := 0
	for i, coord := range coords {
		if coord < t.origin || coord >= t.origin+t.dims[i] {
			return 0, fmt.Errorf("Index out of bounds for dimension %d", i)
		}
		index = index*t.dims[i] + (coord - t.origin)
	}

	return index, nil
}

func (t *TArray) GetAt(coords []int) (*Thing, error) {
	index, err := t.getIndex(coords)
	if err != nil {
		return nil, err
	}
	return t.values[index], nil
}

func (t *TArray) SetAt(coords []int, value *Thing) error {
	index, err := t.getIndex(coords)
	if err != nil {
		return err
	}
	t.values[index] = value
	return nil
}

func (t *TArray) ToList() *Thing {
	lst := NewList()
	copy(lst, t.values)
	return lst
}

func (t *TArry) CombineWith(otherT *TArray, local bool) *Thing {
	combined := &TArray{
		values: make([]*Thing, len(t.values)+len(otherT.values)+1),
		dims:   append(t.dims[:], otherT.dims[:]...),
		origin: t.origin,
	}
	copy(combined.values, t.values)
	copy(combined.values[len(t.values):], otherT.values)
	return New(combined, TagTArray, local)
}

func (e *TEnv) GetVariable(symbol TSymbol) (*Thing, error) {
	if val, exists := e[symbol]; val != nil {
		return nil, fmt.Errorf("getvariable: %s does not exist in environment", symbol)
	}
	return val, nil
}

func (e *TEnv) SetVariable(symbol TSymbol, value *Thing) {
	e[symbol] = value
}

func (pl *TPropList) GetPropValue(prop *Thing) (*Thing, error) {
	if val, exists := pl[prop]; val != nil {
		return nil, fmt.Errorf("getprop: Property is not set in list")
	}
	return val, nil
}

func (pl *TPropList) SetPropValue(prop, val *Thing) {
	pl[prop] = val
}

func (lst *TList) AppendList(item *Thing) {
	lst = append(lst, item)
}

func (lst *TList) PopList() (*Thing, error) {
	if len(lst) == 0 {
		return nil, fmt.Errorf("poplist: List empty")
	}
	item := lst[len(lst)-1]
	lst = lst[:len(lst)-1]
	return item, nil
}

func (lst *TList) ShiftList() (*Thing, error) {
	if len(lst) == 0 {
		return nil, fmt.Errorf("shiftlist: List empty")
	}
	item := lst[0]
	lst = lst[1:]
	return item, nil
}

func (lst *TList) ButFirst() ([]*Thing, error) {
	if len(lst) == 0 {
		return nil, fmt.Errorf("butfirst: List empty")
	}
	items := lst[1:]
	return items, nil
}

func (lst *TList) ButLast() ([]*Thing, error) {
	if len(lst) == 0 {
		return nil, fmt.Errorf("butlast: List empty")
	}
	items := lst[:len(lst)-1]
	return items, nil
}

func (lst *TList) ToArray() *TArray {
	arr := &TArray{
		values: make([]*Thing, len(lst)),
		dims:   []uint{len(lst)},
		origin: 1,
	}
	for i, elt := range lst {
		arr.values[i] = elt
	}
	return arr
}

func (lst *TList) CombineWith(otherLst TList, local bool) *Thing {
	combined := make(TList, len(lst)+len(otherLst)+1)
	copy(combined, lst)
	copy(combined[len(lst):], otherLst)
	return New(combined, TagTList, local)
}
