package main

import (
	"errors"
	"fmt"
	termbox "github.com/nsf/termbox-go"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// returns true if she loves me, false if she loves me not
func randomBool() bool {
	// Intn returns [0,n)
	return rand.Intn(2) == 0
}

/*
makes a new field with ships pseudo-randomly dispersed in the dumbest way
possible
TODO: needs serious work
*/
func makeScatteredField(d dimensions, generics []ship, owner player) field {
	ships := make([]ship, len(generics))
	for shipI, s := range generics {
		var horizontal bool
		var c coord
		var lastc coord
		var coords []coord
		var occupied bool
	TryCoords:
		for {
			c = d.randomCoord()
			horizontal = randomBool()

			if horizontal {
				lastc = c.right(s.length - 1)
			} else {
				lastc = c.down(s.length - 1)
			}

			if !lastc.within(d) {
				continue TryCoords
			}

			// start determining the coordinates the ship in the water takes up
			coords = make([]coord, s.length)
			for i := 0; i < s.length; i++ {
				if horizontal {
					coords[i] = c.right(i)
					// FIXME: repeated call
					occupied = coordOccupied(c.right(i), ships)
				} else {
					coords[i] = c.down(i)
					// FIXME: repeated call
					occupied = coordOccupied(c.down(i), ships)
				}
				if occupied {
					continue TryCoords
				}
			}
			s.spaces = coords
			break // TODO: make less horrendous
		}
		s.owner = owner
		ships[shipI] = s
	}
	return field{
		dimensions: d,
		ships:      ships,
		misses:     []coord{},
	}
}

// given an int, returns the alphabet letter associated with this number
// 'a': 0, 'b': 1, 'c': 2, ... 'y': 24, 'z': 25
func alphabetPosition(s string) (int, error) {
	letter := []rune(s)[0]
	first := []rune("a")[0]
	if letter < first {
		return 0, errors.New("column (lettered index) out of range!")
	}
	return int(letter - first), nil
}

// given an int, returns the alphabet letter associated with this number
// 0: 'a', 1: 'b', 2: 'c', ... 24: 'y', 25: 'z'
func letterInPosition(n int) string {
	first := []rune("a")[0]
	return string(first + rune(n))
}

// a high-level wrapper to prompt the user for a move and call shoot() on a
// field
func move(f *field) {
	var raw []byte
	fmt.Scanf("%s", &raw)
	rawCoord := string(raw)
	rowLetter := strings.ToLower(string(rawCoord[0:1]))
	row, err := alphabetPosition(rowLetter)
	if err != nil {
		panic(err)
	}
	column, err := strconv.Atoi(string(rawCoord[1:]))
	if err != nil {
		panic(err)
	}
	aim := coord{row, column}
	hit, hitShip := f.shoot(aim)
	if hit {
		fmt.Println("You hit a", hitShip.name)
	} else {
		fmt.Println("Miss")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	termbox.Init()
	defer termbox.Close()

	dim := dimensions{10, 10}

	ships := canonicalBattleship()

	attackField := makeScatteredField(dim, ships, adversary)
	defendField := makeScatteredField(dim, ships, human)

	for attackField.shipsLeft() && defendField.shipsLeft() {
		defendField.Show()
		move(&attackField)
	}

	if defendField.shipsLeft() {
		fmt.Println("You've won! Congratulations.")
	} else if attackField.shipsLeft() {
		fmt.Println("You've lost!")
	}
}
