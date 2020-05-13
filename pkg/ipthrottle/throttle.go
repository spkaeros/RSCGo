/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package ipThrottle

import (
	"time"

	"github.com/spkaeros/rscgo/pkg/strutil"
)

type ipThrottle map[int][]time.Time

func (l ipThrottle) Add(ip string) {
	l[strutil.IPToInteger(ip)] = append(l[strutil.IPToInteger(ip)], time.Now())
}

//Recent returns the number of entries that match the provided IP which were added within the past specified timeFrame
func (l ipThrottle) Recent(ip string, timeFrame time.Duration) int {
	valid := 0
	var removing []string
	if attempts, ok := l[strutil.IPToInteger(ip)]; ok {
		for _, v := range attempts {
			if time.Since(v) < timeFrame {
				valid++
				continue
			}
			removing = append(removing, ip)
		}
		for i := range removing {
			if i == len(attempts) {
				l[strutil.IPToInteger(ip)] = attempts[:i]
			} else {
				l[strutil.IPToInteger(ip)] = append(attempts[:i], attempts[i+1:]...)
			}
		}
	}
	return valid
}

type NetworkThrottle interface {
	Recent(string, time.Duration) int
	Add(string)
}

func NewThrottle() NetworkThrottle {
	return make(ipThrottle)
}
