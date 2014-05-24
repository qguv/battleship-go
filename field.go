package main

import (
	"fmt"
	"reflect"
	"math/rand"
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
func (f field) shoot(aim coord) (bool, *ship) {
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

func (f field) rows() []rune {
	rows := f.dimensions.y
	labels := make([]rune, rows)

	for r := 0; r < rows; r++ {
		// indexes first (and only) element of string to give a rune
		charrune := rune(letterInPosition(r)[0])
		labels[r] = charrune
	}

	return labels
}

func (f field) cols() []rune {
	cols := f.dimensions.x
	labels := make([]rune, cols)

	for c := 0; c < cols; c++ { // not a pun I swear
		// indexes first (and only) element of string to give a rune
		numrune := rune(strconv.Itoa(c)[0])
		labels[c] = numrune
	}

	return labels
}

func (f field) Show() {
	cols := f.dimensions.x
	rows := f.dimensions.y
	ships := f.ships

	fmt.Print(" ")
	for i := 0; i < cols; i++ {
		fmt.Print(" ", i)
	}

	fmt.Print("\n")

	for r := 0; r < rows; r++ {
		fmt.Print(letterInPosition(r))
		for c := 0; c < cols; c++ { // not a pun I swear
			if coordOccupied(coord{c, r}, ships) {
				fmt.Print(" +")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("\n")
	}
}
