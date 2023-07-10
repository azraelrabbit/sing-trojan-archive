package buffer

import (
	"github.com/sagernet/sing/common/atomic"
	F "github.com/sagernet/sing/common/format"
)

type Chunk struct {
	data       []byte
	references atomic.Int64
}

func NewChunk(size int) *Chunk {
	var chunk *Chunk
	if size > MaxChunkSize {
		chunk = &Chunk{data: make([]byte, size)}
	} else {
		chunk = getChunkPool(size).Get().(*Chunk)
		for i := range chunk.data {
			chunk.data[i] = 0
		}
	}
	chunk.references.Store(1)
	return chunk
}

//go:nosplit
func (c *Chunk) IncRef() {
	v := c.references.Add(1)
	if v <= 1 {
		panic(F.ToString("Incrementing non-positive count ", v, " on buffer"))
	}
}

//go:nosplit
func (c *Chunk) TryIncRef() bool {
	const speculativeRef = 1 << 32
	if v := c.references.Add(speculativeRef); int32(v) == 0 {
		c.references.Add(-speculativeRef)
		return false
	}
	c.references.Add(-speculativeRef + 1)
	return true
}

func (c *Chunk) ReadRefs() int64 {
	return c.references.Load()
}

//go:nosplit
func (c *Chunk) DecRef() {
	v := c.references.Add(-1)
	switch {
	case v < 0:
		panic(F.ToString("Decrementing non-positive ref count ", v, " on buffer"))
	case v == 0:
		if len(c.data) < MaxChunkSize {
			getChunkPool(len(c.data)).Put(c)
		}
	}
}

func (c *Chunk) Clone() *Chunk {
	newChunk := NewChunk(len(c.data))
	copy(newChunk.data, c.data)
	return newChunk
}
