/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-16-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package isaac

type isaacError struct {
	msg string
}

func (err *isaacError) Error() string {
	return err.msg
}
