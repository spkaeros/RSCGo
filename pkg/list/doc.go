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
package list

/**
 * Package list implements a generic slice-backed list API that can hold any types, and the types don't even have
 * to be uniform, though in practice they probably should be.  I may phase out this code in the near future,
 * as honestly it is nothing more than a wrapper around a standard go slice.  I am not sure, but I believe it is
 * VERY similar to the container/list package in the Go standard library.
 */
