package buf2

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/sagernet/sing/common"
)

func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.Size() == 0 {
		return 0, io.EOF
	}
	n := copy(p, b.AsSlice())
	b.Advance(n)
	return n, nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if b == nil {
		panic("cannot write to a nil view")
	}
	if len(p) == 0 {
		return
	}
	if b.Full() {
		return 0, io.ErrShortBuffer
	}
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	n = copy(b.chunk.data[b.end:], p)
	b.end += n
	return
}

func (b *Buffer) ReadAt(p []byte, off int) (int, error) {
	if off < 0 || off > b.Size() {
		return 0, fmt.Errorf("ReadAt(): offset out of bounds: want 0 < off < %d, got off=%d", b.Size(), off)
	}
	n := copy(p, b.AsSlice()[off:])
	return n, nil
}

func (b *Buffer) WriteAt(p []byte, off int) (int, error) {
	if b == nil {
		panic("cannot write to a nil view")
	}
	if off < 0 || off > b.Size() {
		return 0, fmt.Errorf("write offset out of bounds: want 0 < off < %d, got off=%d", b.Size(), off)
	}
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	n := copy(b.AsSlice()[off:], p)
	if n < len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.Size() == 0 {
		return 0, io.EOF
	}
	p := b.AsSlice()[0]
	b.start++
	return p, nil
}

func (b *Buffer) WriteByte(d byte) error {
	if b == nil {
		panic("cannot write to a nil view")
	} else if b.Full() {
		return io.ErrShortBuffer
	}
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	b.chunk.data[b.end] = d
	b.end++
	return nil
}

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	if b.Size() > 0 {
		sz := b.Size()
		m, e := w.Write(b.AsSlice())
		b.Advance(m)
		n = int64(m)
		if e != nil {
			return n, e
		}
		if m != sz {
			return n, io.ErrShortWrite
		}
	}
	return n, nil
}

func (b *Buffer) ReadFrom0(r io.Reader) (n int64, err error) {
	if b == nil {
		panic("cannot write to a nil view")
	}
	if b.AvailableSize() == 0 {
		return 0, io.ErrShortBuffer
	}
	nInt, err := r.Read(b.availableSlice())
	b.end += nInt
	n += int64(nInt)
	if err == io.EOF {
		return n, nil
	}
	return
}

func (b *Buffer) WriteRandom(size int) []byte {
	buffer := b.Extend(size)
	common.Must1(io.ReadFull(rand.Reader, buffer))
	return buffer
}

func (b *Buffer) sharesChunk() bool {
	return b.chunk.references.Load() > 1
}

func (b *Buffer) availableSlice() []byte {
	if b.sharesChunk() {
		defer b.chunk.DecRef()
		b.chunk = b.chunk.Clone()
	}
	return b.chunk.data[b.end:]
}
