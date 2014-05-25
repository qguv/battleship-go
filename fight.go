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
func move(f *field) string {
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
	hit, hitShip := (&f).shoot(aim)
	if hit {
		return fmt.Sprintln("You hit a", hitShip.name)
	} else {
		return fmt.Sprintln("Miss")
	}
}

func buildLabels(f field, origin coord) {
	var labelOrigin coord

	// column headers
	labelOrigin.x = origin.x
	labelOrigin.y = origin.y - 1
	for i, symbol := range f.cols() {
		theoretical := coord{i, 0}
		actual := theoretical.viewPos(labelOrigin)
		termbox.SetCell(actual.x, actual.y, symbol, termbox.ColorMagenta, termbox.ColorBlack)
	}

	// row headers
	labelOrigin.x = origin.x - 2
	labelOrigin.y = origin.y
	for i, symbol := range f.rows() {
		theoretical := coord{0, i}
		actual := theoretical.viewPos(labelOrigin)
		termbox.SetCell(actual.x, actual.y, symbol, termbox.ColorMagenta, termbox.ColorBlack)
	}
}

func buildInnerField(f field, origin coord) {
	var symbol rune
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack

	for y := 0; y < f.dimensions.y; y++ {
		for x := 0; x < f.dimensions.x; x++ {
			theoretical := coord{x, y}
			statusHere := f.statusAt(theoretical)

			switch statusHere {
			case unknown:
				symbol = ' '
			case empty:
				symbol = '~'
				fg = termbox.ColorCyan
			case miss:
				symbol = '~'
				fg = termbox.ColorGreen
			case hit:
				symbol = '▓'
				fg = termbox.ColorRed
			case occupied:
				symbol = '█'
			}

			actual := theoretical.viewPos(origin)
			termbox.SetCell(actual.x, actual.y, symbol, fg, bg)
		}
	}
}


func counterMove(f *field) string {
	var aim coord
	for {
		aim = f.dimensions.randomCoord()
		for _, c := range f.misses {
			if c == aim {
				continue
			}
		}
		for _, s := range f.ships {
			if s.statusAt(aim) == hit {
				continue
			}
		}
		hit, hitShip := (&f).shoot(aim)
		if hit {
			var destroyed string
			punctuation := "."
			if hitShip.isDestroyed() {
				destroyed = " and destroyed"
				punctuation = "!"
			}
			return fmt.Sprintf("Adversary hit%v your %v%v", destroyed, hitShip.name, punctuation)
		}
	}
}

func tbprint(origin coord, msg string, fg, bg termbox.Attribute) coord {
	for _, r := range msg {
		termbox.SetCell(origin.x, origin.y, r, fg, bg)
		origin = origin.right(1)
	}
	return origin
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dim := dimensions{10, 10}

	ships := canonicalBattleship()

	attackField := makeScatteredField(dim, ships, adversary)
	attackOrigin := coord{0, 0}
	defendField := makeScatteredField(dim, ships, human)
	defendOrigin := coord{30, 0}

	printOrigin := coord{0, 30}

	termbox.Init()
	defer termbox.Close()
	termbox.HideCursor()
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	buildLabels(attackField, attackOrigin)
	buildLabels(defendField, defendOrigin)

	// game loop
	var msg_us, msg_them string
	for attackField.shipsLeft() && defendField.shipsLeft() {
		buildInnerField(attackField, attackOrigin)
		buildInnerField(defendField, defendOrigin)
		termbox.Flush()
		msg_us = move(&attackField)
		termbox.Flush()
		msg_them = counterMove(&defendField)
		tbprint(printOrigin, msg_us, termbox.ColorWhite, termbox.ColorBlack)
		tbprint(printOrigin.down(1), msg_them, termbox.ColorWhite, termbox.ColorBlack)
		termbox.Flush()
	}

	if defendField.shipsLeft() {
		fmt.Println("You've won! Congratulations.")
	} else if attackField.shipsLeft() {
		fmt.Println("You've lost!")
	}
}
