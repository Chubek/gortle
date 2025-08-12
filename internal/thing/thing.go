package thing

import (
	"fmt"
	"strings"
)

type TArray struct {
	Values []Thing
	Dims   uint
	Size   uint
}

type TProc struct {
	Env     TEnv
	NParams int
	Body    *Thing
}

type TEnv map[TSymbol]Thing

type TString string
type TNumber float64
type TSymbol string

type Thing struct {
	Value interface{}
	Next  *Thing
	Tail  *Thing
}

func NewArrayThing(values []Thing, dims, size uint) *TArray {
	arr := &TArray{
		Values: values,
		Dims:   dims,
		Size:   size,
	}
	return arr
}

func NewProcThing(env TEnv, nparams uint, body *Thing) *TProc {
	proc := &TProc{
		Env:     env,
		NParams: nparams,
		Body:    body,
	}
	return proc
}

func New(value interface{}) {
	tng := &Thing{
		Value: value,
		Next:  nil,
		Tail:  nil,
	}
	tng.Tail = tng
	return tng
}

func (tng *Thing) AppendValue(value interface{}) {
	tng.Tail.Next = New(value)
	tng.Tail = tng.Tail.Next
}
