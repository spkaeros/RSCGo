/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package isaac

/**
 * Package isaac contains code that implements a cryptographically secure pseudorandom number generator designed
 * by Bob J. Jenkins Jr. in 1993.  It has been examined by at least a handful of academics in an effort to
 * discover any weaknesses in the results it produces, e.g any patterns, any biases, any way to discover
 * the key used from the results generated, and it has stood up to scrutiny over the years, despite some
 * who see how simple the implementation is and assume that it must not be as good as more sophisticated
 * algorithms.  Some biases were found in a number of certain keys, however, the academic who discovered these
 * biases, Jean-Philippe Aumasson, also proposed a few modifications to fix these biases, as well as some other
 * modifications to improve the algorithm in other ways.  I have read his report about it, and implemented
 * the modifications he recommended, even though even if I hadn't, the biases hardly affect the security
 * of ISAAC.  Practically speaking, ISAAC is accepted by the cryptographic community as providing 256-bit security
 * and I can't see how using it would be any less secure than using ChaCha20, Salsa20, or AES-256.  It also runs
 * in software and is orders of magnitude faster in software than AES, and honestly I haven't got a clue if
 * ChaCha/Salsa variants are slower or faster than ISAAC.
 *
 * In summary, this package implements ISAAC with Jean-Philippe Aumasson's proposed modifications to avoid bias.
 */
