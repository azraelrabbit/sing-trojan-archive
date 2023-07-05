package buf2

import (
	"strconv"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &Buffer{}
	},
}

type Buffer struct {
	start int
	end   int
	chunk *Chunk
}

func New(capacity int) *Buffer {
	chunk := NewChunk(capacity)
	buffer := bufferPool.Get().(*Buffer)
	*buffer = Buffer{chunk: chunk}
	return buffer
}

func From(data []byte) *Buffer {
	buffer := New(len(data))
	buffer.Write(data)
	return buffer
}

func (b *Buffer) Cut(start int, end int) *Buffer {
	if b == nil {
		panic("cannot cut a nil view")
	}
	b.chunk.IncRef()
	newBuffer := bufferPool.Get().(*Buffer)
	newBuffer.chunk = b.chunk
	newBuffer.start = start
	newBuffer.end = end
	return newBuffer
}

func (b *Buffer) Release() {
	if b == nil {
		panic("cannot release a nil view")
	}
	b.chunk.DecRef()
	*b = Buffer{}
	bufferPool.Put(b)
}

func (b *Buffer) Reset() {
	if b == nil {
		panic("cannot reset a nil view")
	}
	b.start = 0
	b.end = 0
}

func (b *Buffer) Full() bool {
	return b == nil || b.end == len(b.chunk.data)
}

func (b *Buffer) Clone() *Buffer {
	if b == nil {
		panic("cannot clone a nil view")
	}
	b.chunk.IncRef()
	newBuffer := bufferPool.Get().(*Buffer)
	newBuffer.chunk = b.chunk
	newBuffer.start = b.start
	newBuffer.end = b.end
	return newBuffer
}

func (b *Buffer) Capacity() int {
	if b == nil {
		return 0
	}
	return len(b.chunk.data)
}

func (b *Buffer) Size() int {
	if b == nil {
		return 0
	}
	return b.end - b.start
}

func (b *Buffer) AsSlice() []byte {
	if b.Size() == 0 {
		return nil
	}
	return b.chunk.data[b.start:b.end]
}

func (b *Buffer) ToSlice() []byte {
	if b.Size() == 0 {
		return nil
	}
	s := make([]byte, b.Size())
	copy(s, b.AsSlice())
	return s
}

func (b *Buffer) AvailableSize() int {
	if b == nil {
		return 0
	}
	return len(b.chunk.data) - b.end
}

func (b *Buffer) AvailableSlice() []byte {
	if b == nil {
		panic("cannot get available slice from a nil view")
	}
	return b.availableSlice()
}

func (b *Buffer) ExtendHeader(n int) []byte {
	if b == nil {
		panic("cannot extend a nil view")
	}
	if b.start < n {
		panic("buffer overflow: cap " + strconv.Itoa(b.Capacity()) + ",start " + strconv.Itoa(b.start) + ", need " + strconv.Itoa(n))
	}
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	b.start -= n
	return b.chunk.data[b.start : b.start+n]
}

func (b *Buffer) Extend(n int) []byte {
	if b == nil {
		panic("cannot extend a nil view")
	}
	newEnd := b.end + n
	if newEnd > b.Capacity() {
		panic("buffer overflow: cap " + strconv.Itoa(b.Capacity()) + ",end " + strconv.Itoa(b.end) + ", need " + strconv.Itoa(n))
	}
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	newSlice := b.chunk.data[b.end : b.end+n]
	b.end = newEnd
	return newSlice
}

func (b *Buffer) Advance(from int) {
	if b == nil {
		panic("cannot advance a nil view")
	}
	b.start += from
}

func (b *Buffer) Truncate(to int) {
	if b == nil {
		panic("cannot truncate a nil view")
	}
	b.end = b.start + to
}
