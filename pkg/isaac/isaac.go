package isaac

import (
	//	"math/bits"
	"sync"

	"github.com/spkaeros/rscgo/pkg/log"
)

// MixMask is used to shave off the first 2 LSB from certain values during result generation
const MixMask = 0xFF << 2

//ISAAC The state representation of the ISAAC CSPRNG
type ISAAC struct {
	// external results
	randrsl [256]uint32
	randcnt int

	// internal state
	state               [256]uint32
	acc1, acc2, counter uint32
	index               int
	remainder           []byte
	sync.RWMutex
}

// this will attempt to shake up the initial state a bit.  Algorithm based off of Mersenne Twister's init
// Not sure if it's of any benefit at all, as ISAAC even when initialized to all zeros is non-uniform
func (r *ISAAC) Seed(seed int64) {
	r.randrsl = padSeed(uint32(seed))
	r.randInit()
}

// ISAAC64+ result shaker, with modifications recommended by Jean-Phillipe Aumasson to avoid some bias,
// to strengthen the output stream.  Replaces an addition with a xor, adds another xor, and
// The below bit rotation ops were changed from right bitshifts of the same number of bits,
// which purportedly gets more diffusion out of existing state bits.
func (r *ISAAC) shake(i uint32, mixed uint32) {
	prevState := r.state[i]
	r.acc1 = mixed + r.state[(i+128)&0xFF]
	// accumulators XOR op changed from ADD op in ISAAC64
	// supposed to remove the reported biases
	//r.state[i] = (r.acc1 ^ r.acc2) + r.state[int(bits.RotateLeft64(prevState&MixMask, 3))&0xFF]
	// ISAAC64 non-modified:
	r.state[i] = (r.acc1 + r.state[(prevState >> 2)&0xFF] + r.acc2)
	// XOR op was added to result value calculation in ISAAC64
	// supposed to reduce the linearity over ZsubText(pow(2,32))
	//	r.acc2 = prevState + (r.acc1 ^ r.state[bits.RotateLeft64(r.state[i]&MixMask, 11)&0xFF])
	// ISAAC64 non-modified:
	r.acc2 = prevState + (r.state[(r.state[i]>>10)&0xFF])
	r.randrsl[i] = r.acc2
}

func (r *ISAAC) generateNextSet() {
	// count
	r.counter++
	// accumulate
	r.acc2 += r.counter

	for i := uint32(0); i < 256; i += 1 {
		// ISAAC64 plus cipher code, with modifications recommended by Jean-Phillipe Aumasson to avoid a discovered bias,
		// and strengthen the output stream.
		r.shake(i, r.acc1^r.acc1<<13)
		i += 1
		r.shake(i, r.acc1^r.acc1>>6)
		i += 1
		r.shake(i, r.acc1^r.acc1<<2)
		i += 1
		r.shake(i, r.acc1^r.acc1>>16)
		/*		switch i % 4 {
				case 0:
					r.acc1 = ^(r.acc1 ^ r.acc1<<21)
				case 1:
					r.acc1 = r.acc1 ^ r.acc1>>5
				case 2:
					r.acc1 = r.acc1 ^ r.acc1<<12
				case 3:
					r.acc1 = r.acc1 ^ r.acc1>>33
				}
				//indirect lookup into the opposite half of state and add to first accumulator
				r.acc1 += r.state[(i+128)&0xFF]
				//store previous state[i] for accumulation with the second accumulator
				prevState := r.state[i]
				// use old state[i] as the basis for an indirect lookup into the state array
				// we mask off the first 3 bits and then shift them off the value,
				// and use that as our indirect access point, we then add that to both
				// accumulators xored with each other to replace that state variable we just stashed.
				r.state[i] = (r.acc1 ^ r.acc2) + r.state[prevState&MixMask>>3]
				// then we use that new state value with a byte shifted off of the start,
				// those pesky first 3 bits that we don't want get masked and shifted off again
				// to get our indirect access point, xor that with the first accumulator, and
				// add it to the previous state variable we stashed earlier, to put some fresh
				// entropy into our second accumulator variable
				r.acc2 = prevState + (r.acc1 ^ r.state[r.state[i]>>8&MixMask>>3])

				// The next acc2 start point is the same as our next result, which is handy
				// since we're now done
				r.randrsl[i] = r.acc2

				 * Original ISAAC cipher code below
				x := r.state[i]
				r.acc1 += r.state[(i+128)&0xFF]           // indirection, accumulation
				y := r.state[(x>>2)&0xFF] + r.acc1 + r.acc2 // indirection, addition, shifts
				r.state[i] = y
				r.acc2 = r.state[(y>>10)&0xFF] + x // indirection, addition, shifts
				r.randrsl[i] = r.acc2
		*/
	}
}

func (r *ISAAC) randInit() {
	const gold = 0x9E3779B9
	var ia [8]uint32//{ gold, gold, gold, gold, gold, gold, gold, gold }
	for i := range ia {
		ia[i] = gold
	}
// 
	// mix1 := func(i int, v, v2 int) {
		// ia[i] ^= v//ia[(i+4)&7]
		// ia[(i+7)&7] += ia[i]
		// ia[(i+4)&7] += v2
	// }
    mix1 := func(i int, v uint32) {
        ia[i] ^= v
        ia[(i+3)%8] += ia[i]
        ia[(i+1)%8] += ia[(i+2)%8]
    }
    mix := func() {
        mix1(0, ia[1]<<11)
        mix1(1, ia[2]>>2)
        mix1(2, ia[3]<<8)
        mix1(3, ia[4]>>16)
        mix1(4, ia[5]<<10)
        mix1(5, ia[6]>>4)
        mix1(6, ia[7]<<8)
        mix1(7, ia[0]>>9)
    }
    // mix := func() {
		// ia[0] ^= ia[4] << 11; ia[7] += ia[0]; ia[4] += ia[1];
		// ia[4] ^= ia[1] >> 2; ia[1] += ia[7]; ia[3] += ia[4];
		// ia[1] ^= ia[7] << 8; ia[2] += ia[1]; ia[7] += ia[3];
		// ia[7] ^= ia[3] >> 16; ia[6] += ia[7]; ia[3] += ia[2];
		// ia[3] ^= ia[2] << 10; ia[2] += ia[6]; ia[5] += ia[3];
		// ia[2] ^= ia[6] >> 4; ia[6] += ia[5]; ia[0] += ia[2];
		// ia[6] ^= ia[5] << 8; ia[4] += ia[6]; ia[5] += ia[0];
		// ia[5] ^= ia[0] >> 9; ia[0] += ia[4]; ia[1] += ia[5];
		// mix1(0, ia[4]<<11, ia[4]+ia[1])
		// mix1(4, ia[1]>>2, ia[3]+ia[4])
		// mix1(1, ia[7]<<8, ia[7]+ia[3])
		// mix1(7, ia[3]>>16, ia[3]+ia[2])
		// mix1(3, ia[2]<<10, ia[2])
		// mix1(2, ia[6]>>4)
		// mix1(6, ia[5]<<8)
		// mix1(5, ia[0]>>9)
	// }
	for i := 0; i < 4; i++ {
		mix()
	}
	messify := func(ia2 [256]uint32) {
		for i := 0; i < 256; i += 8 { // fill state[] with messy stuff
			for i1, v := range ia2[i : i+8] {
				ia[i1] += v
			}
			// ia[0] += ia2[i + 0]
			// ia[1] += ia2[i + 2]
			// ia[2] += ia2[i + 5]
			// ia[3] += ia2[i + 4]
			// ia[4] += ia2[i + 1]
			// ia[5] += ia2[i + 7]
 			// ia[6] += ia2[i + 6]
			// ia[7] += ia2[i + 3]
			mix()
			// r.state[i+0] = ia[0]
			// r.state[i+1] = ia[4]
			// r.state[i+2] = ia[1]
			// r.state[i+3] = ia[7]
			// r.state[i+4] = ia[3]
			// r.state[i+5] = ia[2]
			// r.state[i+6] = ia[6]
			// r.state[i+7] = ia[5]
			for i1, v := range ia {
				r.state[i+i1] = v
			}
		}
	}
	r.Lock()
	messify(r.randrsl)
	messify(r.state)

	r.generateNextSet() /* fill in the first set of results */
	r.randcnt = 255       /* reset the counter for the first set of results */
	r.Unlock()
}

//Uint64 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint64() (number uint64) {
	return uint64(r.Uint32() << 32 | r.Uint32())
}

//Int63 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
// Guarenteed non-negative
func (r *ISAAC) Int63() (number int64) {
	return (int64(r.Uint32()) << 1 >> 1) << 32 | int64(r.Uint32())
}

//Uint32 Returns the next 4 bytes as an integer from the ISAAC CSPRNG receiver instance.
// Guarenteed non-negative
func (r *ISAAC) Uint32() (number uint32) {
	r.Lock()
	defer r.Unlock()
	number = uint32(r.randrsl[r.randcnt])
	r.randcnt--
	if r.randcnt == 0 {
		r.generateNextSet()
		r.randcnt = 255
	}
	return
//	return uint32(r.Int63() >> 31)
}

//Int31 returns a non-negative pseudo-random 31-bit integer as an int32
func (r *ISAAC) Int31() int32 {
	return int32(r.Uint32() << 1 >> 1)
}

func (r *ISAAC) Int() int {
	return int(r.Int31())
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

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// Returns -1 if provided upper-bound <= 0
func (r *ISAAC) Int63n(n int64) int64 {
	if n <= 0 {
		return -1
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
	return r.Uint8n(95) + 0x20
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
	r.Lock()
	defer r.Unlock()
	buf := make([]byte, n)
	for i := 0; i < n;  {
		if len(r.remainder) > 0 && i < len(r.remainder) {
			buf[i] = r.remainder[i]
			r.remainder = r.remainder[i:]
			i += 1
			continue
		}
		r.Unlock()
		nextInt := r.Uint32()
		r.Lock()
		for j := 0; j < 4; j++ {
			next := byte(nextInt >> uint((3-j)<<3))
			if i+j >= n {
				r.remainder = append(r.remainder, next)
				continue
			}
			buf[i+j] = next
		}
		i += 4
	}

	return buf
}

func (r *ISAAC) Float64() float64 {
again:
	f := float64(r.Int63n(1<<53)) / (1 << 53)
	if f == 1 {
		goto again
	}
	return f
}

func (r *ISAAC) Float32() float32 {
again:
	f := float32(r.Float64())
	if f == 1 {
		goto again
	}
	return f
}

// padSeed returns a 256-entry uint64 array filled with values that have been mutated to provide a better initial state.
// Initial padding algorithm copied out of an implementation of the Mersenne twister.
func padSeed(key ...uint32) (seed [256]uint32) {
	if len(key) > 256 {
		log.Warn("Problem initializing ISAAC64+ PRNG seed: Provided key too long; only 256 values will be used.")
	}

	for i := range seed {
		if i < len(key) {
			seed[i] = key[i]
			continue
		}
		seed[i] = 0
	}
	return
}

//New Returns a new ISAAC CSPRNG instance.
func New(key ...uint32) *ISAAC {
	stream := &ISAAC{randrsl: padSeed(key...)}
	stream.randInit()
	return stream
}
