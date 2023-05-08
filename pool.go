package app

// A Pool of reusable items.
type Pool[T any] struct {
	free       chan T
	next       func() T
	factory    func() T
	reset      func(T)
	recycle    func(T)
	configured bool
}

// Next returns the next free item from the pool. This will block unless the
// WithPoolFactory option has been used.
func (p *Pool[T]) Next() (v T) {
	return p.next()
}

// Recycle tries to add the item to the pool. If the WithPoolDiscard option has been
// used then this function will not block if the pool is full.
func (p *Pool[T]) Recycle(v T) {
	if p.reset != nil {
		p.reset(v)
	}
	p.recycle(v)
}

func (p *Pool[T]) blockOnNext() T {
	return <-p.free
}

func (p *Pool[T]) makeOnNext() (v T) {
	select {
	case v = <-p.free:
	default:
		v = p.factory()
	}
	return
}

func (p *Pool[T]) blockOnRecycle(v T) {
	p.free <- v
}

func (p *Pool[T]) discardOnRecycle(v T) {
	select {
	case p.free <- v:
	default:
	}
}

// NewPool returns a pool of the given size having the given options. A pool
// created without options:
//   - blocks on Next() until the pool is no longer empty
//   - does not reset items returned to the pool with Recycle()
//   - blocks on Recycle() until the pool is no longer full.
//
// Pools created without a factory option need to be filled by calling Recycle
// as many times as necessary. When the pool has a factory option the pool is
// filled automatically by the NewPool function.
func NewPool[T any](size int, options ...func(*Pool[T])) (pool *Pool[T]) {
	pool = &Pool[T]{
		free: make(chan T, size),
	}
	pool.next = pool.blockOnNext
	pool.recycle = pool.blockOnRecycle
	for _, opt := range options {
		opt(pool)
	}
	pool.configured = true
	if pool.factory == nil {
		return
	}
	for i := 0; i < size; i++ {
		pool.free <- pool.factory()
	}
	return
}

// WithPoolFactory returns an option to use a factory when Next() is called and
// the pool is empty. The given function must return an initialised item ready
// to use.
func WithPoolFactory[T any](fn func() T) func(*Pool[T]) {
	return func(pool *Pool[T]) {
		if pool.configured {
			return
		}
		pool.factory = fn
		pool.next = pool.makeOnNext
	}
}

// WithPoolReset returns an option to use a function to reset the item before it is
// returned to the pool.
func WithPoolReset[T any](fn func(T)) func(*Pool[T]) {
	return func(pool *Pool[T]) {
		if pool.configured {
			return
		}
		pool.reset = fn
	}
}

// WithPoolDiscard returns an option to discard an item when it is to be recycled and
// the pool is full.
func WithPoolDiscard[T any]() func(*Pool[T]) {
	return func(pool *Pool[T]) {
		if pool.configured {
			return
		}
		pool.recycle = pool.discardOnRecycle
	}
}
