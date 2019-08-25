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

package isaac

type ISAAC struct {
	// external results
	randrsl [256]uint64
	randcnt uint64

	// internal state
	mm         [256]uint64
	aa, bb, cc uint64
	index      int
	remainder  []byte
}

func (r *ISAAC) generateNextSet() {
	r.cc++       // count
	r.bb += r.cc // accumulation

	for i := 0; i < 256; i++ {
		x := r.mm[i]
		// shift
		switch i % 4 {
		case 0:
			// complement, supposedly causes poor states to become randomized quicker.  Avalanche effect.
			r.aa = ^(r.aa ^ r.aa << 21)
		case 1:
			r.aa ^= r.aa >> 5
		case 2:
			r.aa ^= r.aa << 12
		case 3:
			r.aa ^= r.aa >> 33
		}
		// ISAAC(p) cipher code, with modifications recommended by Jean-Phillipe Aumasson to avoid a discovered bias,
		// and strengthen the result set produced
		r.aa += r.mm[(i+128)&0xFF]                       // indirection, accumulation
		y := r.mm[((x>>2)|(x<<30))&0xFF] + (r.aa ^ r.bb) // indirection, addition, (p) exlusive-or, (p) rotation
		r.mm[i] = y
		r.bb = r.aa ^ r.mm[((y>>10)|(y<<22))&0xFF] + x // indirection, addition, (p) exlusive-or, (p) rotation
		r.randrsl[i] = r.bb

		// Original ISAAC cipher code
/*		r.aa += r.mm[(i+128)&0xFF]           // indirection, accumulation
		y := r.mm[(x>>2)&0xFF] + r.aa + r.bb // indirection, addition, shifts
		r.mm[i] = y
		r.bb = r.mm[(y>>10)&0xFF] + x // indirection, addition, shifts
		r.randrsl[i] = r.bb*/
	}
}

/* if (flag==true), then use the contents of randrsl[] to initialize mm[]. */
func (r *ISAAC) randInit() {
	const gold = 0x9e3779b97f4a7c13
	ia := [8]uint64{gold, gold, gold, gold, gold, gold, gold, gold}

	mix1 := func(i int, v uint64) {
		ia[i] -= ia[(i+4)%8]
		ia[(i+5)%8] ^= v
		ia[(i+7)%8] += ia[i]
	}
	mix := func() {
		mix1(0, ia[7]>>9)
		mix1(1, ia[0]<<9)
		mix1(2, ia[1]>>23)
		mix1(3, ia[2]<<15)
		mix1(4, ia[3]>>14)
		mix1(5, ia[4]<<20)
		mix1(6, ia[5]>>17)
		mix1(7, ia[6]<<14)
	}
	for i := 0; i < 4; i++ {
		mix()
	}
	messify := func(ia2 [256]uint64) {
		for i := 0; i < 256; i += 8 { // fill mm[] with messy stuff
			for i1, v := range ia2[i : i+8] {
				ia[i1] += v
			}
			mix()
			for i1, v := range ia {
				r.mm[i+i1] = v
			}
		}
	}
	messify(r.randrsl)
	messify(r.mm)

	r.generateNextSet() /* fill in the first set of results */
	r.randcnt = 0       /* reset the counter for the first set of results */
}

//Seed I might remove this.  I seed exactly once per instance, and I really need at least 4 times this much entropy.
func (r *ISAAC) Seed(key int64) {
	var rsl [256]uint32
	for i := 0; i < 256; i += 2 {
		rsl[i] = uint32(key >> 32)
		rsl[i+1] = uint32(key)
	}
	r.randInit()
}

//Uint64 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint64() (number uint64) {
	number = r.randrsl[r.randcnt]
	r.randcnt++
	if r.randcnt == 256 {
		r.generateNextSet()
		r.randcnt = 0
	}
	return
}

//Uint32 Returns the next 4 bytes as an integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint32() (number uint32) {
	return uint32(r.Uint64())
}

//Int63 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Int63() (number int64) {
	return int64(r.Uint32())<<32 | int64(r.Uint32())
}

//Intn Returns the next 4 bytes as a signed integer of at least 32 bits, with an upper bound of n from the ISAAC CSPRNG.
func (r *ISAAC) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}
	if n <= 1<<31-1 {
		return int(r.Int31n(int32(n)))
	}
	return int(r.Int63n(int64(n)))
}

//Int31n Returns the next 4 bytes as a signed integer of 32 bits, with an upper bound of n from the ISAAC CSPRNG.
func (r *ISAAC) Int31n(n int32) int32 {
	v := r.Uint32()
	prod := uint64(v) * uint64(n)
	low := uint32(prod)
	if low < uint32(n) {
		thresh := uint32(-n) % uint32(n)
		for low < thresh {
			v = r.Uint32()
			prod = uint64(v) * uint64(n)
			low = uint32(prod)
		}
	}
	return int32(prod >> 32)
}

// Int returns a non-negative pseudo-random int.
func (r *ISAAC) Int() int {
	u := uint(r.Int63())
	return int(u << 1 >> 1) // clear sign bit if int == int32
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (r *ISAAC) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return r.Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := r.Int63()
	for v > max {
		v = r.Int63()
	}
	return v % n
}

//Uint8n Returns the next byte as an unsigned 8-bit integer, with an upper bound of n from the ISAAC CSPRNG.
func (r *ISAAC) Uint8n(bound byte) (number byte) {
	for number = r.Uint8(); number < 0 || number >= bound; number = r.Uint8() {
	}
	return
}

//Uint16 Returns the next 2 bytes as a short integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint16() uint16 {
	buf := r.NextBytes(2)
	return uint16(buf[0])<<8 | uint16(buf[1])
}

//Uint8 Returns the next byte from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint8() byte {
	return r.NextBytes(1)[0]
}

//NextChar Returns the next ASCII character from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) NextChar() byte {
	return byte(r.Int31n(95)) + 32
}

//String Returns the next `len` ASCII characters from the ISAAC CSPRNG receiver instance as a Go string.
func (r *ISAAC) String(len int) (ret string) {
	for i := 0; i < len; i++ {
		ret += string(r.NextChar())
	}

	return
}

func (r *ISAAC) Read(dst []byte) (n int, err error) {
	if len(dst) > 0 {
		n = len(dst)
		for i, b := range r.NextBytes(n) {
			dst[i] = b
		}
		return
	}
	return 0, &isaacError{msg: "isaac.Read([]byte): Length of `dst` buffer must be >= 0"}
}

//NextBytes Returns the next `n` bytes from the ISAAC CSPRNG receiver instance, and since ISAAC generates 4-byte words,
//  if you request a length of bytes that is not divisible evenly by 4, it will stash the remaining bytes into a buffer
//  to be used on your next call to this function.
func (r *ISAAC) NextBytes(n int) []byte {
	buf := make([]byte, n)
	r.index = 0
	if len(r.remainder) > 0 {
		for i := 0; i < len(r.remainder) && r.index < n; i++ {
			buf[r.index] = r.remainder[i]
			r.index++
		}
		if r.index >= n {
			r.remainder = r.remainder[r.index:]
			return buf
		}
	}
	r.remainder = []byte{}

	for r.index < n {
		nextInt := r.Uint32()
		for i := 0; i < 4; i++ {
			if r.index >= n {
				r.remainder = append(r.remainder, byte(nextInt>>uint(8*(3-i))))
				continue
			}
			buf[r.index] = byte(nextInt >> uint(8*(3-i)))
			r.index++
		}
	}

	return buf
}

//New Returns a new ISAAC CSPRNG instance.
func New(key []uint64) *ISAAC {
	var tmpRsl [256]uint64
	for i := 0; i < len(key); i++ {
		tmpRsl[i] = key[i]
	}
	if len(key) < 256 {
		for i := len(key); i < 256; i++ {
			tmpRsl[i] = 3735928559 + uint64(i)
		}
	}
	stream := &ISAAC{randrsl: tmpRsl}
	stream.randInit()
	return stream
}
