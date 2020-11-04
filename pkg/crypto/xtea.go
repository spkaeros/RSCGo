package crypto

import "encoding/binary"

const phi = 0x9E3779B9

type XteaKeys struct {
	Keys []int
}

//Decrypt Takes an XTEA block as input, and attempts to decrypt it using
// the stored keys.
func (x *XteaKeys) Decrypt(in []byte) []byte {
	out := make([]byte, len(in))
	blocks := len(in) >> 3

	i := 0
	for ; i < blocks; i++ {
		word1 := binary.BigEndian.Uint32(in[i<<3:])
		word2 := binary.BigEndian.Uint32(in[i<<3+4:])
		sum := uint64(phi << 5)

		for j := 0; j < 1<<5; j++ {
			word2 -= uint32(uint64(((word1 << 4) ^ (word1 >> 5)) + word1) ^ uint64(sum + uint64(x.Keys[(sum >> 11) & 3])))
			sum -= phi
			word1 -= uint32(uint64(((word2 << 4) ^ (word2 >> 5)) + word2) ^ uint64(sum + uint64(x.Keys[sum & 3])))
		}
		binary.BigEndian.PutUint32(out[i<<3:], word1)
		binary.BigEndian.PutUint32(out[i<<3+4:], word2)
	}
	// RSClassic protocol 235 does not pad out XTEA block's size to a multiple of 8
	// The result is that we must append the raw remaining bytes after the closest multiple of 8
	// This leaks any username less than 6 characters long.  It breaks nothing to pad on the client side,
	// so it's highly suggested to do as much to avoid this behavior.
	for i <<= 3; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}
