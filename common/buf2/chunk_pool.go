package buf2

import (
	"fmt"
	"math/bits"
	"sync"
)

const (
	baseChunkSizeLog2 = 6
	baseChunkSize     = 1 << baseChunkSizeLog2          // 64
	MaxChunkSize      = baseChunkSize << (numPools - 1) // 64k
	numPools          = 11
)

var chunkPools [numPools]sync.Pool

func init() {
	for i := 0; i < numPools; i++ {
		chunkSize := baseChunkSize * (1 << i)
		chunkPools[i].New = func() any {
			return &Chunk{
				data: make([]byte, chunkSize),
			}
		}
	}
}

func getChunkPool(size int) *sync.Pool {
	idx := 0
	if size > baseChunkSize {
		idx = bits.Len32(uint32(size)>>6) - 1
		if size > 1<<(idx+baseChunkSizeLog2) {
			idx++
		}
	}
	if idx >= numPools {
		panic(fmt.Sprintf("pool for chunk size %d does not exist", size))
	}
	return &chunkPools[idx]
}
