package rand

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/spkaeros/rscgo/pkg/isaac"
)

var IsaacRng *isaac.ISAAC

func init() {
	var rsl = make([]uint64, 256)
	if err := binary.Read(rand.Reader, binary.BigEndian, rsl); err != nil {
		fmt.Println("ERROR: Could not read ints fully into init slice.", err)
	}
	IsaacRng = isaac.New(rsl)
}

//Uint8 Gets a single random byte of data from the PRNG
func Uint8() uint8 {
	return IsaacRng.Uint8()
}

//Uint8n Gets a single random byte of data from the ISAAC PRNG, bound being the upper bound
func Uint8n(bound uint8) uint8 {
	return IsaacRng.Uint8n(bound)
}

//Uint16 Gets 2 random bytes of data from the PRNG and returns them as a single 16-bit short integer
func Uint16() uint16 {
	return IsaacRng.Uint16()
}

//Uint32 Gets 4 random bytes of data from the PRNG and returns them as a single 32-bit integer
func Uint32() uint32 {
	return IsaacRng.Uint32()
}

//Int31n Returns a randomized 31-bit signed integer from the ISAAC instance seeded by the system CSPRNG, landing between 0 and bound.
func Int31n(bound int) int {
	return int(IsaacRng.Int31n(int32(bound)))
}

//Int31N Returns a randomized 31-bit signed integer from the ISAAC instance seeded by the system CSPRNG, landing between low and high.
func Int31N(low, high int) int {
	return Int31n((high+1)-low) + low
}

//Uint64 Gets 8 random bytes of data from the PRNG and returns them as a single 64-bit long integer
func Uint64() uint64 {
	return IsaacRng.Uint64()
}

//Int63n Gets 8 random bytes of data from the PRNG and returns them as a single 64-bit long integer
func Int63n(bound int) int64 {
	return IsaacRng.Int63n(int64(bound))
}

//Int63n Gets 8 random bytes of data from the PRNG and returns them as a single 64-bit long integer
func Int63N(low, high int) int64 {
	return IsaacRng.Int63n(int64((high+1)-low)) + int64(low)
}

//String Gets n random bytes of data from the PRNG and returns them as a Go string
func String(length int) string {
	return IsaacRng.String(length)
}
