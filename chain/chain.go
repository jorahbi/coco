package chain

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type ChainInterface interface {
	Apply(fn func() error, warps ...string)
	Error() error
}

type chain struct {
	err error
}

var chainPool = sync.Pool{
	New: func() any {
		return &chain{}
	},
}

func NewChain() *chain {
	c := chainPool.Get().(*chain)
	c.err = nil

	return c
}

func (c *chain) Apply(fn func() error, warps ...string) {
	if c.err != nil {
		return
	}
	c.err = errors.Wrap(fn(), strings.Join(warps, ":"))
}

func (c *chain) Error() error {
	defer c.put()
	return c.err
}

func (c *chain) put() {
	c.err = nil
	chainPool.Put(c)
}
