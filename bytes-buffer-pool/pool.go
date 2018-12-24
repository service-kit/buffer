package bytes_buffer_pool

import (
	"bytes"
	"sync"
)

type Pool struct {
	pool sync.Pool
}

func (p *Pool) init() {
	p.pool.New = func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 256*1<<10))
	}
}

func (p *Pool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

func (p *Pool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

func NewBytesBufferPool() *Pool {
	p := new(Pool)
	p.init()
	return p
}
