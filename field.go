package main

import (
	"errors"
	"math/rand"
	"reflect"
	"strconv"
)

// coord and dimensions both contain x and y fields, but these fields have
// different meaning in each type and therefore the methods they avail are
// distinct.
type coord struct{ x, y int }
type dimensions struct{ x, y int }

func (from coord) dimensions(to coord) dimensions {
	x := to.x - from.x
	y := to.y - from.y
	return dimensions{x, y}
}

// returns the coordinate with an y-value increased by the amount specified
func (c coord) down(y1 int) coord {
	y0 := c.y
	return coord{c.x, y0 + y1}
}

// returns the coordinate with an x-value increased by the amount specified
func (c coord) right(x1 int) coord {
	x0 := c.x
	return coord{x0 + x1, c.y}
}

// returns the area between two coords
func (c0 coord) area(c1 coord) int {
	return c0.dimensions(c1).area()
}

// returns the area contained in a dimension field
func (d dimensions) area() int {
	return d.x * d.y
}

// returns an arbitrary coordinate within the specified dimensions
func (d dimensions) randomCoord() (c coord) {
	c.x = rand.Intn(d.x)
	c.y = rand.Intn(d.y)
	return
}

// coord.within returns true if a coordinate is within the dimensions
// specified, starting at the origin.
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

// returns the player name as a string
func (p player) String() (s string) {
	switch p {
	case nobody:
		s = "nobody"
	case human:
		s = "human"
	case adversary:
		s = "adversary"
	case everybody:
		s = "everybody"
	}
	return s
}

type ship struct {
	name   string
	length int
	spaces []coord
	holes  []coord
	owner  player
}

// tests whether a ship has been fatally shot
func (s ship) isDestroyed() bool {
	return reflect.DeepEqual(s.spaces, s.holes)
}

type field struct {
	dimensions dimensions
	ships      []ship
	misses     []coord
	owner      player
}

// returns whether there are any remaining ships which aren't destroyed
func (f field) shipsLeft() bool {
	for _, s := range f.ships {
		if !s.isDestroyed() {
			return true
		}
	}
	return false
}

// tests whether a coord is a valid coordinate in a field
func (c coord) on(f field) bool {
	return c.within(f.dimensions)
}

// Shoot at a coordinate. If a ship is located at the coordinate, mutate the
// ship to indicate its damage and return the address of the hit ship.
func (f *field) shoot(aim coord) (bool, *ship) {
	for _, s := range f.ships {
		for _, c := range s.spaces {
			if c == aim {
				s.holes = append(s.holes, aim)
				return true, &s
			}
		}
	}
	return false, &ship{}
}

func coordOccupied(aim coord, ships []ship) bool {
	for _, s := range ships {
		for _, c := range s.spaces {
			if c == aim {
				return true
			}
		}
	}
	return false
}

// status is an enumerated type indicating status of specific ship coordinates
type status uint8

const (
	unknown status = iota
	miss
	hit
	empty
	occupied
)

// If the ship isn't here, return miss.
// If the ship is here, return occupied or hit depending on whether already
// damaged.
func (s ship) statusAt(aim coord) status {
	for _, c := range s.spaces {
		if c == aim {
			// either a hit or just a ship here
			for _, h := range s.holes {
				if aim == h {
					return hit
				}
			}
			return occupied
		}
	}
	return miss
}

func (f field) statusAt(aim coord) status {
	// if it's a known miss, show that
	for _, c := range f.misses {
		if c == aim {
			return miss
		}
	}
	var shipOwner player
	for _, s := range f.ships {
		shipOwner = s.owner
		shipStatus := s.statusAt(aim)

		// nothing to see--we move on to the next ship or execute the bit after
		// the loop
		if shipStatus == miss {
			continue
		}

		// we know about our own field, and everyone knows about hits
		if shipOwner == human || shipStatus == hit {
			return shipStatus

			// since empty spaces are taken care of with continue, the only
			// possibility is that it's occupied, and we shouldn't know that
		} else if shipStatus == occupied {
			return unknown

			// if I've made a logical error or overlooked a case...
		} else {
			errors.New("don't know what happened here... check the logic in field.statusAt")
		}
	}

	// here, there is no data for the coordinate
	if f.owner == human {
		// if it's our field, we know it's empty
		return empty
	}

	// if it's not our field, we don't know
	return unknown
}

func (f field) rows() []rune {
	rows := f.dimensions.y
	labels := make([]rune, rows)

	for r := 0; r < rows; r++ {
		character, err := letterInPosition(r)
		if err != nil {
			panic(err)
		}
		// indexes first (and only) element of string to give a rune
		labels[r] = rune(character[0])
	}

	return labels
}

func (f field) cols() []rune {
	cols := f.dimensions.x
	labels := make([]rune, cols)
	over10 := []rune{'⒑', '⒒', '⒓', '⒔', '⒕', '⒖', '⒗', '⒘', '⒙', '⒚', '⒛'}

	for c := 1; c <= cols; c++ { // not a pun I swear
		// indexes first (and only) element of string to give a rune
		if c < 10 {
			numrune := rune(strconv.Itoa(c)[0])
			labels[c-1] = numrune
		} else if c <= 20 {
			labels[c-1] = over10[c-10]
		} else {
			panic(errors.New("not enough monospace numbers to represent dimensions"))
		}
	}

	return labels
}

// This function tells us where in termbox a coordinate would be represented.
func (theoretical coord) viewPos(origin coord) (actual coord) {
	actual.x = origin.x + theoretical.x*2 + 4
	actual.y = origin.y + theoretical.y + 2
	return
}
