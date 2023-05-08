package app

import (
	"testing"
)

func TestWithPoolFactory(t *testing.T) {
	type data struct {
		value int
	}
	pool := NewPool(
		1,
		WithPoolFactory(func() *data { return &data{} }),
	)
	pool.Next()
	x := pool.Next()
	if x == nil {
		t.Error()
	}
}

func TestWithPoolReset(t *testing.T) {
	type data struct {
		value int
	}
	pool := NewPool(
		1,
		WithPoolFactory(func() *data { return &data{} }),
		WithPoolReset(func(d *data) { d.value = 0 }),
	)
	x := pool.Next()
	x.value = 1
	pool.Recycle(x)
	x = pool.Next()
	if x.value != 0 {
		t.Error()
	}
}

func TestWithPoolDiscard(t *testing.T) {
	type data struct {
		value int
	}
	pool := NewPool(
		1,
		WithPoolFactory(func() *data { return &data{value: 1} }),
		WithPoolDiscard[*data](),
	)
	pool.Recycle(&data{value: 2})
	x := pool.Next()
	if x.value != 1 {
		t.Error()
	}
	x = pool.Next()
	if x.value != 1 {
		t.Error()
	}
}
