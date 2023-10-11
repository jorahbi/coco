package chain

import (
	"errors"
	"fmt"
	"testing"
)

var ErrTest = errors.New("test error1")

func TestChain(t *testing.T) {
	c := NewChain()
	c.Apply(func() error {
		return ErrTest
	})
	c.Apply(func() error {
		return fmt.Errorf("error %v", "2")
	})
	err := c.Error()
	fmt.Println(errors.Is(err, ErrTest), errors.Unwrap(err))
}
