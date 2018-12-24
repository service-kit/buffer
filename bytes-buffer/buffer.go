package bytes_buffer

import (
	"errors"
	"math"
)

type BufferNode struct {
	dataEnd   int
	dataBegin int
	cap       int
	size      int
	buffer    []byte
	next      *BufferNode
	pre       *BufferNode
	id        int
}

func (node *BufferNode) Bytes() []byte {
	var out []byte = nil
	if node.dataEnd > node.dataBegin {
		out = node.buffer[node.dataBegin:node.dataEnd]
	} else {
		out = make([]byte, node.size)
		n := copy(out, node.buffer[node.dataBegin:])
		copy(out[n:], node.buffer[:node.dataEnd])
	}
	return out
}

func (node *BufferNode) Write(in []byte) (int, error) {
	if nil == node || nil == node.buffer {
		return 0, errors.New("not init")
	}
	if node.Full() {
		return 0, nil
	}
	inLen := len(in)
	n := copy(node.buffer[node.dataEnd:], in)
	if node.dataBegin > 0 && n < inLen {
		m := copy(node.buffer[:node.dataBegin], in[n:])
		node.dataEnd = m
		n += m
	} else {
		node.dataEnd += n
		node.dataEnd %= node.cap
	}
	node.size += n
	return n, nil
}

func (node *BufferNode) Read(out []byte) (int, error) {
	if nil == node || nil == node.buffer {
		return 0, errors.New("not init")
	}
	if 0 == node.size {
		return 0, nil
	}
	outLen := len(out)
	n := 0
	if node.dataBegin < node.dataEnd {
		n = copy(out, node.buffer[node.dataBegin:node.dataEnd])
	} else {
		n = copy(out, node.buffer[node.dataBegin:])
		if outLen > n {
			n += copy(out[n:], node.buffer[:node.dataEnd])
		}
	}
	node.dataBegin += n
	node.size -= n
	if 0 == node.size {
		node.clear()
	}
	return n, nil
}

func (node *BufferNode) clear() {
	node.size = 0
	node.dataBegin = 0
	node.dataEnd = 0
}

func (node *BufferNode) ReadAll() []byte {
	if nil == node || nil == node.buffer {
		return nil
	}
	if 0 == node.size {
		return nil
	}
	var out []byte = nil
	if node.dataEnd > node.dataBegin {
		out = node.buffer[node.dataBegin:node.dataEnd]
	} else {
		out = make([]byte, node.size)
		n := copy(out, node.buffer[node.dataBegin:])
		copy(out[n:], node.buffer[:node.dataEnd])
	}
	node.clear()
	return out
}

func (node *BufferNode) Full() bool {
	return node.cap == node.size
}

func (node *BufferNode) Len() int {
	return node.size
}

type BytesBuffer struct {
	root     *BufferNode
	writePos *BufferNode
	readPos  *BufferNode
	size     int
	cap      int
	count    int
	perCap   int
}

func (buf *BytesBuffer) init(cap int) {
	buf.perCap = 1 << 16
	buf.cap = buf.perCap
	buf.size = 0
	buf.count = 1
	buf.root = NewNode(buf.perCap)
	buf.root.next = buf.root
	buf.root.pre = buf.root
	buf.writePos = buf.root
	buf.readPos = buf.root
	buf.add(cap - 1)
}

func (buf *BytesBuffer) Cap() int {
	return buf.cap
}

func (buf *BytesBuffer) Len() int {
	return buf.size
}

func (buf *BytesBuffer) Write(in []byte) (int, error) {
	inLen := len(in)
	if buf.cap < buf.size+inLen {
		buf.ensureCapacity(inLen - (buf.cap - buf.size))
	}
	inIdx := 0
	for {
		if inIdx >= inLen {
			break
		}
		if nil == buf.writePos {
			return 0, errors.New("buffer not space")
		}
		n, err := buf.writePos.Write(in[inIdx:])
		if nil != err {
			return 0, err
		}
		inIdx += n
		buf.size += n
		if buf.writePos.Full() {
			buf.writePos = buf.writePos.next
		}
	}
	return inLen, nil
}

func (buf *BytesBuffer) Read(out []byte) (int, error) {
	outLen := len(out)
	outIdx := 0
	for {
		if outIdx >= outLen {
			break
		}
		if nil == buf.readPos || (buf.writePos == buf.readPos && 0 == buf.readPos.Len()) {
			break
		}
		n, err := buf.readPos.Read(out[outIdx:])
		if nil != err {
			return 0, err
		}
		outIdx += n
		buf.size -= n
		if 0 == buf.readPos.Len() {
			buf.readPos = buf.readPos.next
		}
	}
	return outLen, nil
}

func (buf *BytesBuffer) Bytes() []byte {
	out := make([]byte, buf.size)
	outIndex := 0
	for iter := buf.readPos; outIndex == buf.size; iter = iter.next {
		outIndex += copy(out[outIndex:], iter.Bytes())
	}
	return out
}

func (buf *BytesBuffer) ensureCapacity(n int) {
	lval := math.Log2(float64((n / buf.perCap) + buf.count))
	lval = math.Ceil(lval)
	newNodeCount := int(math.Pow(2, lval)) - buf.count
	for i := 0; i < newNodeCount; i++ {
		newNode := NewNode(buf.perCap)
		if 0 == i && buf.writePos == buf.readPos {
			newNode.Write(buf.writePos.ReadAll())
			buf.writePos = newNode
		}
		buf.insertNode(newNode, buf.readPos)
	}
}

func (buf *BytesBuffer) add(n int) {
	for i := 0; i < n; i++ {
		newNode := NewNode(buf.perCap)
		buf.insertNode(newNode, buf.root)
	}
}

func (buf *BytesBuffer) insertNode(newNode, pos *BufferNode) {
	pos.pre.next = newNode
	newNode.pre = pos.pre
	newNode.next = pos
	pos.pre = newNode
	newNode.id = buf.count
	buf.count++
	buf.cap += buf.perCap
}

func NewNode(cap int) *BufferNode {
	node := new(BufferNode)
	node.cap = cap
	node.size = 0
	node.dataEnd = 0
	node.buffer = make([]byte, cap)
	return node
}

func NewBytesBuffer(cap int) *BytesBuffer {
	buffer := new(BytesBuffer)
	buffer.init(cap)
	return buffer
}
