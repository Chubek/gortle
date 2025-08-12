package thing

import (
	"fmt"
	"strings"

	"gortle/internal/runtime"
)

const growthFactorList = 0.75

type TArray struct {
	Values []*Thing
	Size   uint
	Origin uint
}

type TProc struct {
	Env     TEnv
	NParams int
	Body    *Thing
}

type TList struct {
	Vals  []*Thing
	Size  uint
	Count uint
}

type TEnv map[TSymbol]Thing

type TString string
type TNumber float64
type TSymbol string

type Thing struct {
	Value interface{}
}

func NewArrayThing(values []*Thing, size, origin uint) *TArray {
	arr := &TArray{
		Values: values,
		Size:   size,
		Origin: origin,
	}
	return arr
}

func (arr *TArray) Index(idx, dim uint) *Thing {
	if idx >= (arr.Size*dim)+arr.Origin {
		runtime.Errorf()
	}
	return arr.Values[(idx*dim)+arr.Origin]
}

func NewProcThing(env TEnv, nparams uint, body *Thing) *TProc {
	proc := &TProc{
		Env:     env,
		NParams: nparams,
		Body:    body,
	}
	return proc
}

func NewListThing(head *Thing, size uint) *TList {
	lst := &TList{
		Vals:  make([]*Thing, size),
		Count: 1,
		Size:  size,
	}
	return lst
}

func (lst *TList) Append(ntng *Thing) {
	if lst.Count/lst.Size >= growthFactorList {
		newLst := make([]*Thing, lst.Size*2)
		copy(newLst, lst.Vals)
		lst.Vals = newLst
		lst.Size *= 2
	}
	lst.Vals[lst.Count] = ntng
	lst.Count += 1
}

func (lst *TList) Index(idx uint) *Thing {
	if idx >= lst.Count {
		runtime.Errorf("Index %d out of bounds", idx)
	}
	return lst.Vals[idx]
}

func (lst *TList) Delete(idx uint) {
	if idx >= lst.Count {
		runtime.Errorf("Index %d out of bounds", idx)
	}
	append(lst.Vals[:idx], lst.Vals[idx+1:]...)
}

func (lst *TList) ButFirst() []*Thing {
	return lst.Vals[1:]
}

func (lst *TList) ButLast() []*Thing {
	return lst.Vals[:len(lst.Vals)-1]
}

func New(value interface{}) *Thing {
	tng := &Thing{
		Value: value,
	}
	return tng
}
