/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-24-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"bitbucket.org/zlacki/rscgo/pkg/isaac"
)

var rscRand *isaac.ISAAC

func init() {
	initRsl := make([]uint32, 256)
	if err := binary.Read(rand.Reader, binary.BigEndian, initRsl); err != nil {
		fmt.Println("ERROR: Could not read ints fully into init slice.", err)
	}
	rscRand = isaac.New(initRsl)
}

//getRandomData Reads n random bytes from the system-specific PRNG and returns them in a byte slice
func getRandomData(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rscRand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

//getCryptoRandomData Reads n random bytes from the system-specific CSPRNG and returns them in a byte slice
func getCryptoRandomData(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

//GetSecureRandomString Gets n random bytes of data from the CSPRNG and returns them as a Go string
func GetSecureRandomString(length int) string {
	b, err := getCryptoRandomData(length)
	if err != nil {
		return "nil"
	}

	return base64.URLEncoding.EncodeToString(b)
}

//GetSecureRandomLong Gets 8 random bytes of data from the CSPRNG and returns them as a single 64-bit long integer
func GetSecureRandomLong() uint64 {
	b, err := getCryptoRandomData(8)
	if err != nil {
		return 0
	}

	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//GetSecureRandomInt Gets 4 random bytes of data from the CSPRNG and returns them as a single 32-bit integer
func GetSecureRandomInt() uint32 {
	b, err := getCryptoRandomData(4)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

//GetSecureRandomSmart Gets 3 random bytes of data from the CSPRNG and returns them as a single 24-bit smart integer
func GetSecureRandomSmart() uint32 {
	b, err := getCryptoRandomData(3)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[2])
}

//GetSecureRandomShort Gets 2 random bytes of data from the CSPRNG and returns them as a single 16-bit short integer
func GetSecureRandomShort() uint16 {
	b, err := getCryptoRandomData(2)
	if err != nil {
		return 0
	}

	return (uint16(b[0]) << 8) | uint16(b[1])
}

//GetSecureRandomByte Gets a single random byte of data from the CSPRNG
func GetSecureRandomByte() uint8 {
	b, err := getCryptoRandomData(1)
	if err != nil {
		return 0
	}

	return uint8(b[0])
}

//GetRandomByte Gets a single random byte of data from the PRNG
func GetRandomByte() uint8 {
	b, err := getRandomData(1)
	if err != nil {
		return 0
	}

	return uint8(b[0])
}

//GetRandomShort Gets 2 random bytes of data from the PRNG and returns them as a single 16-bit short integer
func GetRandomShort() uint16 {
	b, err := getRandomData(2)
	if err != nil {
		return 0
	}

	return (uint16(b[0]) << 8) | uint16(b[1])
}

//GetRandomSmart Gets 3 random bytes of data from the PRNG and returns them as a single 24-bit smart integer
func GetRandomSmart() uint32 {
	b, err := getRandomData(3)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 16) | (uint32(b[1]) << 8) | uint32(b[2])
}

//GetRandomInt Gets 4 random bytes of data from the PRNG and returns them as a single 32-bit integer
func GetRandomInt() uint32 {
	b, err := getRandomData(4)
	if err != nil {
		return 0
	}

	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])
}

//GetRandomLong Gets 8 random bytes of data from the PRNG and returns them as a single 64-bit long integer
func GetRandomLong() uint64 {
	b, err := getRandomData(8)
	if err != nil {
		return 0
	}

	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//GetRandomString Gets n random bytes of data from the PRNG and returns them as a Go string
func GetRandomString(length int) string {
	b, err := getRandomData(length)
	if err != nil {
		return "nil"
	}

	return base64.URLEncoding.EncodeToString(b)
}
