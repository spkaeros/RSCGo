/*
------------------------------------------------------------------------------
resetRandomBuffer.go: an implementation of Bob Jenkins' random number generator ISAAC based on 'readable.c'
* 18 Aug 2014 -- direct port of readable.c to Go
* 10 Sep 2014 -- updated to be more idiomatic Go
------------------------------------------------------------------------------
*/

package isaac

type ISAAC struct {
	// external results
	randrsl [256]uint32
	randcnt uint32

	// internal state
	mm         [256]uint32
	aa, bb, cc uint32

	index int
	remainder []byte
}

func (r *ISAAC) resetRandomBuffer() {
	r.cc++ // cc gets incremented once per 256 results
	r.bb += r.cc // then combined with bb

	for i := 0; i < 256; i++ {
		x := r.mm[i]
		switch i % 4 {
		case 0:
			r.aa ^= r.aa << 13
		case 1:
			r.aa ^= r.aa >> 6
		case 2:
			r.aa ^= r.aa << 2
		case 3:
			r.aa ^= r.aa >> 16
		}
		r.aa += r.mm[(i+128)%256]
		y := r.mm[(x>>2)%256] + r.aa + r.bb
		r.mm[i] = y
		r.bb = r.mm[(y>>10)%256] + x
		r.randrsl[i] = r.bb
	}

	/* Note that bits 2..9 are chosen from x but 10..17 are chosen
	   from y.  The only important thing here is that 2..9 and 10..17
	   don't overlap.  2..9 and 10..17 were then chosen for speed in
	   the optimized version (rand.c) */
	/* See http://burtleburtle.net/bob/rand/isaac.html
	   for further explanations and analysis. */
}

func mix(a, b, c, d, e, f, g, h uint32) (uint32, uint32, uint32, uint32, uint32, uint32, uint32, uint32) {
	a ^= b << 11
	d += a
	b += c
	b ^= c >> 2
	e += b
	c += d
	c ^= d << 8
	f += c
	d += e
	d ^= e >> 16
	g += d
	e += f
	e ^= f << 10
	h += e
	f += g
	f ^= g >> 4
	a += f
	g += h
	g ^= h << 8
	b += g
	h += a
	h ^= a >> 9
	c += h
	a += b
	return a, b, c, d, e, f, g, h
}

/* if (flag==true), then use the contents of randrsl[] to initialize mm[]. */
func (r *ISAAC) randInit(flag bool) {
	var a, b, c, d, e, f, g, h uint32
	a, b, c, d, e, f, g, h = 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9

	for i := 0; i < 4; i++ {
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
	}

	for i := 0; i < 256; i += 8 { /* fill mm[] with messy stuff */
		if flag { /* use all the information in the seed */
			a += r.randrsl[i]
			b += r.randrsl[i+1]
			c += r.randrsl[i+2]
			d += r.randrsl[i+3]
			e += r.randrsl[i+4]
			f += r.randrsl[i+5]
			g += r.randrsl[i+6]
			h += r.randrsl[i+7]
		}
		a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
		r.mm[i] = a
		r.mm[i+1] = b
		r.mm[i+2] = c
		r.mm[i+3] = d
		r.mm[i+4] = e
		r.mm[i+5] = f
		r.mm[i+6] = g
		r.mm[i+7] = h
	}

	if flag { /* do a second pass to make all of the seed affect all of mm */
		for i := 0; i < 256; i += 8 {
			a += r.mm[i]
			b += r.mm[i+1]
			c += r.mm[i+2]
			d += r.mm[i+3]
			e += r.mm[i+4]
			f += r.mm[i+5]
			g += r.mm[i+6]
			h += r.mm[i+7]
			a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
			r.mm[i] = a
			r.mm[i+1] = b
			r.mm[i+2] = c
			r.mm[i+3] = d
			r.mm[i+4] = e
			r.mm[i+5] = f
			r.mm[i+6] = g
			r.mm[i+7] = h
		}
	}

	r.resetRandomBuffer() /* fill in the first set of results */
	r.randcnt = 256       /* reset the counter for the first set of results */
}

//Seed This method has been changed so as to reflect the Jagex style seeding the client uses.
//  Takes a total of 16 bytes of entropy, usually half comes from the server cryptographically-secure PRNG, and the
//  other half from the client's cryptographically secure PRNG.
func (r *ISAAC) Seed(key []uint32) {
	for i, k := range key {
		r.randrsl[i] = k
	}
	r.randInit(true)
}

//Uint64 Returns the next 8 bytes as a long integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint64() uint64 {
	return uint64(r.Uint32()) << 32 | uint64(r.Uint32())
}

//Uint32 Returns the next 4 bytes as an integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint32() (number uint32) {
	r.randcnt--
	number = r.randrsl[r.randcnt]
	if r.randcnt == 0 {
		r.resetRandomBuffer()
		r.randcnt = 256
	}
	return
}

func (r *ISAAC) Int31n(n int) int32 {
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

func (r *ISAAC) Uint8n(bound byte) (number byte) {
	for number = r.Uint8(); number < 0 || number >= bound; number = r.Uint8() {  }
	return
}

//Uint16 Returns the next 2 bytes as a short integer from the ISAAC CSPRNG receiver instance.
func (r *ISAAC) Uint16() uint16 {
	buf := r.NextBytes(2)
	return uint16(buf[0]) << 8 | uint16(buf[1])
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
		} else {
			r.remainder = []byte{}
		}
	}

	for ; r.index < n; {
		nextInt := r.Uint32()
		if n % 4 != 0 && n - r.index < 4 {
			spaceLeft := n - r.index
			for i := 0; i < spaceLeft; i++ {
				buf[r.index] = byte(nextInt >> uint(8*(3-i)))
				r.index++
			}
			r.remainder = []byte{}
			for i := spaceLeft; i < 4; i++ {
				r.remainder = append(r.remainder, byte(nextInt >> uint(8*(3-i))))
			}
		} else {
			for i := 0; i < 4; i++ {
				buf[r.index] = byte(nextInt >> uint(8*(3-i)))
				r.index++
			}
		}
	}

	return buf
}

func NewISAACStream(key []uint32) *ISAAC {
	stream := new(ISAAC)
	stream.Seed(key)
	return stream
}
