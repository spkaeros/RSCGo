package crypto

func DecryptXtea(data []byte, offset, length int, keys []int) []byte {
	blocks := length >> 3
	ret := make([]byte, length)

	in := data[offset:offset+length]

	i := 0
	for ; i < blocks; i++ {
		v0, v1 := blockToInt(in[i<<3:])
		sum := uint64(0x9E3779B9 << 5)

		for j := 0; j < 1<<5; j++ {
			v1 -= uint32(uint64(((v0 << 4) ^ (v0 >> 5)) + v0) ^ uint64(sum + uint64(keys[(sum >> 11) & 3])))
			sum -= 0x9E3779B9
			v0 -= uint32(uint64(((v1 << 4) ^ (v1 >> 5)) + v1) ^ uint64(sum + uint64(keys[sum & 3])))
		}
		intToBlock(v0, v1, ret[i<<3:])
	}
	for i <<= 3; i < len(data); i++ {
		ret[i] = data[i]
	}
	return ret
}

// blockToUint32 reads an 8 byte slice into two uint32s.
// The block is treated as big endian.
func blockToInt(src []byte) (uint32, uint32) {
	r0 := uint32(src[0])<<24 | uint32(src[1])<<16 | uint32(src[2])<<8 | uint32(src[3])
	r1 := uint32(src[4])<<24 | uint32(src[5])<<16 | uint32(src[6])<<8 | uint32(src[7])
	return r0, r1
}

// uint32ToBlock writes two uint32s into an 8 byte data block.
// Values are written as big endian.
func intToBlock(v0, v1 uint32, dst []byte) {
	dst[0] = byte(v0 >> 24)
	dst[1] = byte(v0 >> 16)
	dst[2] = byte(v0 >> 8)
	dst[3] = byte(v0)
	dst[4] = byte(v1 >> 24)
	dst[5] = byte(v1 >> 16)
	dst[6] = byte(v1 >> 8)
	dst[7] = byte(v1)
}
