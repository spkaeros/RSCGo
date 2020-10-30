package rand

import (
	stdrand "math/rand"
	"time"
	"github.com/spkaeros/rscgo/pkg/isaac"

	// "github.com/spkaeros/rscgo/pkg/log"
	// "github.com/spkaeros/rscgo/pkg/strutil"
)

var fastLockedSrc *isaac.ISAAC

var Rng *stdrand.Rand

func init() {
	// key := []byte("This is my secret key")
	// key := []byte{1}
	// i := 0
	// var iv = make([]uint32, len(key)>>2)
	// for i = 0; i>>2 < len(key)>>2; i += 4 {
		// j := i>>2
		// i1 := uint32(key[i]) << 0 | uint32(key[i+1]) << 8 | uint32(key[i+2]) << 16 | uint32(key[i+3]) << 24
		// iv[j] = i1
	// 
	// }

	fastLockedSrc = isaac.New(uint32(time.Now().UnixNano()))
	// for _,v := range iv {
		// log.Debug(v, strutil.Base16.String(uint64(v)))
	// }
	// fastLockedSrc = isaac.New(1)
	// for i := 0; i < 256; i++ {
		// v := fastLockedSrc.Uint32()
		// 
		// log.Debug(v, strutil.Base16.String(uint64(v)))
	// }
	Rng = stdrand.New(fastLockedSrc)
	stdrand.Seed(Rng.Int63())
}

func Source() *isaac.ISAAC {
	return fastLockedSrc
}
