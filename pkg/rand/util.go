package rand

import (
	rand2 "math/rand"
	"time"
	"github.com/spkaeros/rscgo/pkg/isaac"

	"github.com/spkaeros/rscgo/pkg/log"
	"github.com/spkaeros/rscgo/pkg/strutil"
)

var fastLockedSrc *isaac.ISAAC

var Rng *rand2.Rand

func init() {
	// arr := []byte("This is my secret key")
	// i := 0
	// for i = 0; i>>2 < len(arr)>>2; i += 4 {
		// j := i>>2
		// i1 := int(arr[i]) << 0 | int(arr[i+1]) << 8 | int(arr[i+2]) << 16 | int(arr[i+3]) >> 24
		// key[j] = i1
	// 
	// }
	// for i := 0; i < len(arr); i++ {
		// key[i] |= uint32(arr[i])
		// log.Debug(strutil.Base16.String(uint64(arr[i])))
	// }

	fastLockedSrc = isaac.New(uint32(time.Now().UnixNano()))
	for i := 0; i < 256; i++ {
		v := fastLockedSrc.Uint32()
		
		log.Debug(v, strutil.Base16.String(uint64(v)))
	}
	// var out = make([]uint64, len("a top secret secret"))//[]uint64{uint64(fastLockedSrc.NextChar()),uint64(fastLockedSrc.NextChar()),uint64(fastLockedSrc.NextChar()),uint64(fastLockedSrc.NextChar()),uint64(fastLockedSrc.NextChar())}
	// for i, v := range []byte("a Top Secret secret") {
		// out[i] = uint64(v^byte(fastLockedSrc.Uint8()%95) + 0x20)
	// }
	// for i := range out {
		// log.Debug(strutil.Base16.String(out[i]))
		// log.Debug(out[i])
// 
    	// log.Debug(fmt.Sprintf("%X", out))
   	// }
	Rng = rand2.New(fastLockedSrc)
	rand2.Seed(Rng.Int63())
}

func Source() *isaac.ISAAC {
	return fastLockedSrc
}
