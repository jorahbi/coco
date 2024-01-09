package op

import "sync"

type ObjectItemInterface interface {
	Reset()
}

func NewPool(o ObjectItemInterface) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return o
		},
	}
}
