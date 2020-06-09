/*
 * Copyright (c) 2020 Zachariah Knight <aeros.storkpk@gmai.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in a copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package entity

type Type int

const (
	TypePlayer     Type = 1 << iota // 1
	TypeNpc                         // 2
	TypeObject                      // 4
	TypeDoor                        // 8
	TypeItem                        // 16
	TypeGroundItem                  // 32

	TypeMob    = TypePlayer | TypeNpc
	TypeEntity = TypeObject | TypeDoor | TypeItem | TypeGroundItem
)

type Entity interface {
	Location
	ServerIndex() int
	Type() Type
}

type Location interface {
	X() int
	Y() int
	Point() Location
}
