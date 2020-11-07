package rand

import (
	stdrand "math/rand"
	"encoding/binary"
	"math/bits"
	crand "crypto/rand"
	// "time"
	"github.com/spkaeros/rscgo/pkg/isaac"

	// "github.com/spkaeros/rscgo/pkg/log"
	// "github.com/spkaeros/rscgo/pkg/strutil"
)

var fastLockedSrc *isaac.ISAAC

var Rng *isaac.ISAAC

func init() {
	// get 1KiB of semantically/cryptographically secure PRNG data
	seed := make([]byte, 1024)
	crand.Read(seed)
	seedWords := make([]int, len(seed)>>2)
	for i := range seedWords {
		seedWords[i] = int(binary.LittleEndian.Uint32(seed[i<<2:]))
	}
	fastLockedSrc = isaac.New(seedWords...)
	// for _,v := range iv {
		// log.Debug(v, strutil.Base16.String(uint64(v)))
	// }
	// fastLockedSrc = isaac.New(1)
	// for i := 0; i < 256; i++ {
		// v := fastLockedSrc.Uint32()
		// 
		// log.Debug(v, strutil.Base16.String(uint64(v)))
	// }
	Rng = fastLockedSrc
	stdrand.Seed(Rng.Int63())
}

func Source() *isaac.ISAAC {
	return fastLockedSrc
}

func Int() int {
	if bits.UintSize <= 32 {
		return int(fastLockedSrc.Int31())
	}
	return int(fastLockedSrc.Int63())
}

func Intn(n int) int {
	if n <= 1<<31-1 {
		return int(fastLockedSrc.Int31n(int32(n)))
	}
	return int(fastLockedSrc.Int63n(int64(n)))
}

func Uintn(n uint) uint {
	if n <= 1<<31-1 {
		return uint(fastLockedSrc.Uint32()%uint32(n))
	}
	return uint(fastLockedSrc.Uint64()%uint64(n))
}

func Uint() uint {
	if bits.UintSize <= 32 {
		return uint(fastLockedSrc.Uint32())
	}
	return uint(fastLockedSrc.Uint64())
}

func Float64() float64 {
	return fastLockedSrc.Float64()
}

func Float32() float32 {
	return fastLockedSrc.Float32()
}

func Byte() byte {
	return fastLockedSrc.Uint8()
}

func Bytes(n int) []byte {
	return fastLockedSrc.NextBytes(n)
}

func String(n int) []byte {
	return []byte(fastLockedSrc.NextBytes(n))
}
