/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-21-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package entity

//Object Represents a game object in the world.
type Object struct {
	ID        int
	Direction int
	Boundary  bool
	location  *Location
	Index     int
}

func (o *Object) X() int {
	return o.location.X
}

func (o *Object) Y() int {
	return o.location.Y
}

func (o *Object) Location() *Location {
	return o.location
}

func NewObject(id, direction, x, y int, boundary bool) *Object {
	return &Object{id, direction, boundary, &Location{x, y}, -1}
}
