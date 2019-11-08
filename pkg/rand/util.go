package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/spkaeros/rscgo/pkg/isaac"
)

var rscRand *isaac.ISAAC

func init() {
	initRsl := make([]uint64, 256)
	if err := binary.Read(rand.Reader, binary.BigEndian, initRsl); err != nil {
		fmt.Println("ERROR: Could not read ints fully into init slice.", err)
	}
	rscRand = isaac.New(initRsl)
}

//RandomBytes Reads n random bytes from the system-specific PRNG and returns them in a byte slice
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rscRand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

//RandomBytesS Reads n random bytes from the system-specific CSPRNG and returns them in a byte slice
func RandomBytesS(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

//StringS Gets n random bytes of data from the CSPRNG and returns them as a Go string
func StringS(length int) string {
	b, err := RandomBytesS(length)
	if err != nil {
		return "nil"
	}

	return base64.URLEncoding.EncodeToString(b)
}

//Uint64S Gets 8 random bytes of data from the CSPRNG and returns them as a single 64-bit long integer
func Uint64S() uint64 {
	b, err := RandomBytesS(8)
	if err != nil {
		return 0
	}

	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//Uint32S Gets 4 random bytes of data from the CSPRNG and returns them as a single 32-bit integer
func Uint32S() uint32 {
	b, err := RandomBytesS(4)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

//Uint24S Gets 3 random bytes of data from the CSPRNG and returns them as a single 24-bit smart integer
func Uint24S() uint32 {
	b, err := RandomBytesS(3)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[2])
}

//Uint16S Gets 2 random bytes of data from the CSPRNG and returns them as a single 16-bit short integer
func Uint16S() uint16 {
	b, err := RandomBytesS(2)
	if err != nil {
		return 0
	}

	return (uint16(b[0]) << 8) | uint16(b[1])
}

//Uint8S Gets a single random byte of data from the CSPRNG
func Uint8S() uint8 {
	b, err := RandomBytesS(1)
	if err != nil {
		return 0
	}

	return uint8(b[0])
}

//Uint8 Gets a single random byte of data from the PRNG
func Uint8() uint8 {
	b, err := RandomBytes(1)
	if err != nil {
		return 0
	}

	return uint8(b[0])
}

//Uint16 Gets 2 random bytes of data from the PRNG and returns them as a single 16-bit short integer
func Uint16() uint16 {
	b, err := RandomBytes(2)
	if err != nil {
		return 0
	}

	return (uint16(b[0]) << 8) | uint16(b[1])
}

//Uint24 Gets 3 random bytes of data from the PRNG and returns them as a single 24-bit smart integer
func Uint24() uint32 {
	b, err := RandomBytes(3)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[2])
}

//Uint32 Gets 4 random bytes of data from the PRNG and returns them as a single 32-bit integer
func Uint32() uint32 {
	b, err := RandomBytes(4)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

//Int31n Returns a randomized 31-bit signed integer from the ISAAC instance seeded by the system CSPRNG, landing between 0 and bound.
func Int31n(bound int) int {
	return int(rscRand.Int31n(int32(bound)))
}

//Int31N Returns a randomized 31-bit signed integer from the ISAAC instance seeded by the system CSPRNG, landing between low and high.
func Int31N(low, high int) int {
	return Int31n(high-low) + low
}

//Uint64 Gets 8 random bytes of data from the PRNG and returns them as a single 64-bit long integer
func Uint64() uint64 {
	b, err := RandomBytes(8)
	if err != nil {
		return 0
	}

	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//String Gets n random bytes of data from the PRNG and returns them as a Go string
func String(length int) string {
	b, err := RandomBytes(length)
	if err != nil {
		return "nil"
	}

	return base64.URLEncoding.EncodeToString(b)
}
