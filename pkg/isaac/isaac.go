package isaac

import (
	"sync"
)

// MixMask is used to shave off the first 2 LSB from certain values during result generation
const MixMask = 0xFF << 2
const phi = 0x9E3779B9

//ISAAC The state representation of the ISAAC CSPRNG
type ISAAC struct {
	// external results
	randrsl [256]uint32
	randcnt int

	// internal state
	state			   [256]uint32
	acc1, acc2, counter uint32
	index			   int
	remainder		   []byte
	sync.RWMutex
}

// this will attempt to shake up the initial state a bit.  Algorithm based off of Mersenne Twister's init
// Not sure if it's of any benefit at all, as ISAAC even when initialized to all zeros is non-uniform
func (r *ISAAC) Seed(seed int64) {
	r.randrsl = padKeys(int(seed))
	r.randInit()
}

func (r *ISAAC) generateNextSet() {
	// count
	r.counter++
	// accumulate
	r.acc2 += r.counter

	var shifts = [...]uint { 13, 6, 2, 16 }
	for i := uint32(0); i < 256; i += 1 {
		prevState := r.state[i]
		var mixed uint32
		if i&1 == 0 {
			mixed = r.acc1^r.acc1 << shifts[i&3]
		} else {
			mixed = r.acc1^r.acc1 >> shifts[i&3]
		}
		r.acc1 = mixed + r.state[(i+128)&0xFF]
		// accumulators XOR op changed from ADD op in ISAAC64
		// supposed to remove the reported biases
		//r.state[i] = (r.acc1 ^ r.acc2) + r.state[int(bits.RotateLeft64(prevState&MixMask, 3))&0xFF]
		// ISAAC64 non-modified:
		r.state[i] = r.acc1 + r.state[(prevState >> 2)&0xFF] + r.acc2
		// XOR op was added to result value calculation in ISAAC64
		// supposed to reduce the linearity over ZsubText(pow(2,32))
		//	r.acc2 = prevState + (r.acc1 ^ r.state[bits.RotateLeft64(r.state[i]&MixMask, 11)&0xFF])
		// ISAAC64 non-modified:
		r.acc2 = prevState + r.state[(r.state[i]>>10)&0xFF]
		r.randrsl[i] = r.acc2
	}
}

func (r *ISAAC) randInit() {
	var mess = [...]uint32 { phi, phi, phi, phi, phi, phi, phi, phi }
	mix := func() {
		var shifts = [...]uint { 11, 2, 8, 16, 10, 4, 8, 9 }
		for i := 0; i < 8; i++ {
			if i&1 == 0 {
				mess[i] ^= (mess[(i+1)&7] << shifts[i&7])
			} else {
				mess[i] ^= (mess[(i+1)&7] >> shifts[i&7])
			}
			mess[(i+3)&7] += mess[i]
			mess[(i+1)&7] += mess[(i+2)&7]
		}
	}
	fillMess := func(state [256]uint32) {
		for i := 0; i < 256; i += 8 { // fill state or result-set with messy stuff derived of the golden ratio
			for i1, v := range state[i : i+8] {
				mess[i1] += v
			}
			mix()
			for i1, v := range mess {
				r.state[i+i1] = v
			}
		}
	}

	// initialize int array to derive messy state from, using initial value derived of the golden ratio because nothing's up our sleeve
	// Mix up this initial state words array 4 times in a row
	for i := 0; i < 4; i++ {
		mix()
	}
	r.Lock()
	fillMess(r.randrsl)
	fillMess(r.state)

	r.generateNextSet() // run core ISAAC algorithm to fill result-set
	r.randcnt = 256 // set first output word to last element of result-set
	r.Unlock()
}

//Uint64 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint64() (number uint64) {
	return uint64(r.Uint32()) << 32 | uint64(r.Uint32())
}

//Int63 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Int63() (number int64) {
	// return int64(r.Uint64())<<1>>1
	// return (int64(r.Int31()) << 32) | int64(r.Uint32()) & 0x7FFFFFFFFFFFFFFF
	return int64(r.Uint64())
}

//Uint32 Returns the next 4 bytes as an integer from the ISAAC CSPRNG receiver instance.
// Guarenteed non-negative
func (r *ISAAC) Uint32() (number uint32) {
	r.Lock()
	defer r.Unlock()
	if r.randcnt <= 0 {
		r.generateNextSet()
		r.randcnt = 256
	}
	r.randcnt--
	return uint32(r.randrsl[r.randcnt])
}

//Int31 returns a pseudo-random 31-bit integer as an int32
func (r *ISAAC) Int31() int32 {
	return int32(r.Uint32())
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
	max := int64((1 << 63) - 1 - ((1<<63)%uint64(n)))
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
func (r *ISAAC) String(n int) (ret string) {
	for i := 0; i < n; i++ {
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
			i += 1
			r.remainder = r.remainder[i:]
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
	if f == 1.0 {
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

func padKeys(key ...int) [256]uint32 {
	var seed [256]uint32
	for i := range seed {
		if i >= len(key) {
			seed[i] = 0
			continue
		}
		seed[i] = uint32(key[i])
	}
	return seed
}

//New Returns a new ISAAC CSPRNG instance.
func New(key ...int) *ISAAC {
	stream := &ISAAC{randrsl: padKeys(key...)}
	stream.randInit()
	return stream
}

//New Returns a new ISAAC CSPRNG instance.
func New32(key ...uint32) *ISAAC {
	var keys [256]uint32
	for i := range keys {
		if i >= len(key) {
			keys[i] = 0
			continue
		}
		keys[i] = key[i]
	}
	stream := &ISAAC{randrsl: keys}
	stream.randInit()
	return stream
}
