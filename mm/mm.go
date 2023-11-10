package mm

import (
	"unsafe"

	"github.com/heiyeluren/xmm"
)

var menory xmm.XMemory

func init() {
	var err error
	menory, err = new(xmm.Factory).CreateMemory(0.7)
	if err != nil {
		panic("CreateMemory fail ")
	}
}

type Object[T any] struct {
	p    unsafe.Pointer
	data *T
}

func MustObject[T any]() Object[T] {
	size := unsafe.Sizeof(new(T))
	p, err := menory.Alloc(size)
	if err != nil {
		return Object[T]{data: new(T), p: nil}
	}
	return Object[T]{data: (*T)(p), p: p}
}

func (o Object[T]) Get() *T {
	return o.data
}

func (o Object[T]) Free() {
	if o.p == nil {
		return
	}
	menory.Free(uintptr(o.p))
}
