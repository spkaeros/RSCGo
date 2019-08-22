/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-14-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package errors

type NetError struct {
	msg    string
	Ping   bool
	Closed bool
}

func (e NetError) Error() string {
	return e.msg
}

//NewNetworkError Returns a new NetError struct with the specified message.
func NewNetworkError(s string) error {
	return NetError{msg: s, Closed: true}
}

//ConnClosed Error to return when the connection closes normally.
var ConnClosed = NewNetworkError("Connection closed.")

//ConnTimedOut Error to return when the connection is inactive for 10 seconds.
var ConnTimedOut = NewNetworkError("Connection timed out.")

//ConnDeadline Error to return when the connection's deadline for reading data can not be properly set.
var ConnDeadline = NewNetworkError("Problem setting connection deadline.")

//BufferOverflow Error to return when we accidentally try to read from an empty packet.
var BufferOverflow = NewNetworkError("Attempted to read too much data from packet; would have caused buffer overflow.")
