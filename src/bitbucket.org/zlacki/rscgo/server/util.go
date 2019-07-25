package server

import (
	"crypto/rand"
	"encoding/base64"
	nrand "math/rand"
	"net"
	"strings"
	"time"
)

func init() {
	nrand.Seed(time.Now().UnixNano())
}

//getCryptoRandomData Reads n random bytes from the system-specific cryptographically secure PRNG and returns them in a byte slice
func getCryptoRandomData(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

//GetSecureRandomString Gets n random bytes of data from the cryptographically secure PRNG and returns them as a Go string
func GetSecureRandomString(length int) string {
	b, err := getCryptoRandomData(length)
	if err != nil {
		return "nil"
	}

	return base64.URLEncoding.EncodeToString(b)
}

//GetSecureRandomLong Gets 8 random bytes of data from the cryptographically secure PRNG and returns them as a single 64-bit long integer
func GetSecureRandomLong() uint64 {
	b, err := getCryptoRandomData(8)
	if err != nil {
		return 0
	}

	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//GetSecureRandomInt Gets 4 random bytes of data from the cryptographically secure PRNG and returns them as a single 32-bit integer
func GetSecureRandomInt() uint32 {
	b, err := getCryptoRandomData(4)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

//GetSecureRandomSmart Gets 3 random bytes of data from the cryptographically secure PRNG and returns them as a single 24-bit smart integer
func GetSecureRandomSmart() uint32 {
	b, err := getCryptoRandomData(3)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[2])
}

//GetSecureRandomShort Gets 2 random bytes of data from the cryptographically secure PRNG and returns them as a single 16-bit short integer
func GetSecureRandomShort() uint16 {
	b, err := getCryptoRandomData(2)
	if err != nil {
		return 0
	}

	return (uint16(b[0]) << 8) | uint16(b[1])
}

//GetSecureRandomByte Gets a single random byte of data from the cryptographically secure PRNG
func GetSecureRandomByte() uint8 {
	b, err := getCryptoRandomData(1)
	if err != nil {
		return 0
	}

	return uint8(b[0])
}

func getIPFromConn(c net.Conn) string {
	parts := strings.Split(c.RemoteAddr().String(), ":")
	if len(parts) < 1 {
		return "nil"
	}
	return parts[0]
}