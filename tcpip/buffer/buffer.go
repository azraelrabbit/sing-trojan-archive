package buffer

import (
	"net/netip"
	"strconv"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &PacketBuffer{}
	},
}

type PacketBuffer struct {
	Destination netip.AddrPort
	start       int
	end         int
	chunk       *Chunk
}

func New(capacity int) *PacketBuffer {
	chunk := NewChunk(capacity)
	buffer := bufferPool.Get().(*PacketBuffer)
	*buffer = PacketBuffer{chunk: chunk}
	return buffer
}

func From(data []byte) *PacketBuffer {
	buffer := New(len(data))
	buffer.Write(data)
	return buffer
}

func (b *PacketBuffer) Cut(start int, end int) *PacketBuffer {
	if b == nil {
		panic("cannot cut a nil view")
	}
	b.chunk.IncRef()
	newBuffer := bufferPool.Get().(*PacketBuffer)
	newBuffer.chunk = b.chunk
	newBuffer.start = start
	newBuffer.end = end
	return newBuffer
}

func (b *PacketBuffer) Release() {
	if b == nil {
		panic("cannot release a nil view")
	}
	b.chunk.DecRef()
	*b = PacketBuffer{}
	bufferPool.Put(b)
}

func (b *PacketBuffer) Reset() {
	if b == nil {
		panic("cannot reset a nil view")
	}
	b.start = 0
	b.end = 0
}

func (b *PacketBuffer) Full() bool {
	return b == nil || b.end == len(b.chunk.data)
}

func (b *PacketBuffer) Clone() *PacketBuffer {
	if b == nil {
		panic("cannot clone a nil view")
	}
	b.chunk.IncRef()
	newBuffer := bufferPool.Get().(*PacketBuffer)
	newBuffer.chunk = b.chunk
	newBuffer.start = b.start
	newBuffer.end = b.end
	return newBuffer
}

func (b *PacketBuffer) Capacity() int {
	if b == nil {
		return 0
	}
	return len(b.chunk.data)
}

func (b *PacketBuffer) Size() int {
	if b == nil {
		return 0
	}
	return b.end - b.start
}

func (b *PacketBuffer) AsSlice() []byte {
	if b.Size() == 0 {
		return nil
	}
	return b.chunk.data[b.start:b.end]
}

func (b *PacketBuffer) ToSlice() []byte {
	if b.Size() == 0 {
		return nil
	}
	s := make([]byte, b.Size())
	copy(s, b.AsSlice())
	return s
}

func (b *PacketBuffer) AvailableSize() int {
	if b == nil {
		return 0
	}
	return len(b.chunk.data) - b.end
}

func (b *PacketBuffer) AvailableSlice() []byte {
	if b == nil {
		panic("cannot get available slice from a nil view")
	}
	return b.availableSlice()
}

func (b *PacketBuffer) ExtendHeader(n int) []byte {
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

func (b *PacketBuffer) Extend(n int) []byte {
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

func (b *PacketBuffer) Advance(from int) {
	if b == nil {
		panic("cannot advance a nil view")
	}
	b.start += from
}

func (b *PacketBuffer) Truncate(to int) {
	if b == nil {
		panic("cannot truncate a nil view")
	}
	b.end = b.start + to
}
