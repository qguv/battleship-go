package main

import "reflect"

// coord and dimensions both contain x and y fields, but these fields have different meaning in each type and therefore the methods they avail are distinct.
type coord struct{ x, y int }
type dimensions struct{ x, y int }

func (from coord) dimensions(to coord) dimensions {
	x := to.x - from.x
	y := to.y - from.y
	return dimensions{x, y}
}

func (c coord) down(y1 int) coord {
	y0 := c.y
	return coord{c.x, y0 + y1}
}

func (c coord) right(x1 int) coord {
	x0 := c.x
	return coord{x0 + x1, c.y}
}

func (c0 coord) area(c1 coord) int {
	return c0.dimensions(c1).area()
}

func (d dimensions) area() int {
	return d.x * d.y
}

// coord.within returns true if a coordinate is within the dimensions specified, starting at the origin.
func (c coord) within(d dimensions) bool {
	if c.x < 0 || c.y < 0 {
		return false
	}
	if c.x >= d.x || c.y >= d.y {
		return false
	}
	return true
}

// player is an enumerated type indicating ownership for ships
type player uint8

const (
	nobody player = iota
	human
	adversary
	everybody
)

type ship struct {
	name   string
	length int
	spaces []coord
	holes  []coord
	owner  player
}

func (s ship) isDestroyed() bool {
	return reflect.DeepEqual(s.spaces, s.holes)
}

type field struct {
	dimensions dimensions
	ships      []ship
	misses     []coord
}

func (f field) winner() player {
	survivor := nobody
	for _, s := range f.ships {
		// If there is an undestroyed ship on the board, its owner has not lost.
		if !s.isDestroyed() {
			if survivor < s.owner {
				survivor = survivor + s.owner
			}
		}
		// If both opponents are alive, nobody has won.
		if survivor >= everybody {
			return nobody
		}
	}
	return survivor
}

func (c coord) on(f field) bool {
	return c.within(f.dimensions)
}

func (f field) shoot(aim coord) (bool, *ship) {
	for _, s := range f.ships {
		for _, c := range s.spaces {
			if c == aim {
				return true, &s
			}
		}
	}
	return false, &ship{}
}
