/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package throttle

import (
	"time"

	"github.com/spkaeros/rscgo/pkg/strutil"
)

type ipThrottle map[int][]time.Time

func (l ipThrottle) Add(ip string) {
	l[strutil.AddrToInt(ip)] = append(l[strutil.AddrToInt(ip)], time.Now())
}

//CountSince returns the number of entries that match the provided IP which were added within the past specified timeFrame
func (l ipThrottle) CountSince(ip string, timeFrame time.Duration) int {
	valid := 0
	var removing []int
	addrIntegral := strutil.AddrToInt(ip)
	if attempts, ok := l[addrIntegral]; ok {
		for _, v := range attempts {
			if time.Since(v) < timeFrame {
				valid++
				continue
			}
			removing = append(removing, addrIntegral)
		}
		for i := range removing {
			if i == len(attempts) {
				l[addrIntegral] = attempts[:i]
			} else {
				l[addrIntegral] = append(attempts[:i], attempts[i+1:]...)
			}
		}
	}
	return valid
}

//Throttle An API that any structs acting as a throttle must provide.
type Throttle interface {
	CountSince(string, time.Duration) int
	Add(string)
}

func NewThrottle() Throttle {
	return make(ipThrottle)
}
