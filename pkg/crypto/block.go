// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Implementation adapted from Needham and Wheeler's paper:
	http://www.cix.co.uk/~klockstone/xtea.pdf

	A precalculated look up table is used during encryption/decryption for values that are based purely on the key.
*/

package crypto

// XTEA is based on 64 rounds.
// const numRounds = 64
const numRounds = 32
