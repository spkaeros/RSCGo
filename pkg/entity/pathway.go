/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-20-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package entity

//Pathway Represents a path for a mobile entity to traverse across the virtual world.
type Pathway struct {
	StartX, StartY  int
	WaypointsX      []int
	WaypointsY      []int
	CurrentWaypoint int
}

//waypointXoffset Returns the offset for the X coordinate of the specified waypoint.
func (p *Pathway) waypointXoffset(w int) int {
	if w >= len(p.WaypointsX) || w == -1 {
		return 0
	}
	return p.WaypointsX[w]
}

//waypointX Returns the X coordinate of the specified waypoint.
func (p *Pathway) waypointX(w int) int {
	return p.StartX + p.waypointXoffset(w)
}

//waypointYoffset Returns the offset for the Y coordinate of the specified waypoint.
func (p *Pathway) waypointYoffset(w int) int {
	if w >= len(p.WaypointsY) || w == -1 {
		return 0
	}
	return p.WaypointsY[w]
}

//waypointY Returns the Y coordinate of the specified waypoint.
func (p *Pathway) waypointY(w int) int {
	return p.StartY + p.waypointYoffset(w)
}

//Waypoint Returns the locattion of the specified waypoint
func (p *Pathway) Waypoint(w int) *Location {
	return &Location{p.waypointX(w), p.waypointY(w)}
}

//Start Returns the locattion of the start of the path
func (p *Pathway) Start() *Location {
	return &Location{p.StartX, p.StartY}
}

//NextTile Returns the next tile for the mob to move to in the pathway.
func (p *Pathway) NextTile(startX, startY int) *Location {
	destX := p.waypointX(p.CurrentWaypoint)
	destY := p.waypointY(p.CurrentWaypoint)
	newLocation := &Location{destX, destY}
	switch {
	case startX > destX:
		newLocation.X = startX - 1
		break
	case startX < destX:
		newLocation.X = startX + 1
		break
	}
	switch {
	case startY > destY:
		newLocation.Y = startY - 1
		break
	case startY < destY:
		newLocation.Y = startY + 1
		break
	}
	return newLocation
}
