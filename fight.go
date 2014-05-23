package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func randomBool() bool {
	// Intn returns [0,n)
	return rand.Intn(2) == 0
}

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
					occupied = coordOccupied(c.right(i), ships) //FIXME: repeated call
				} else {
					coords[i] = c.down(i)
					occupied = coordOccupied(c.down(i), ships) //FIXME: repeated call
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

func alphabetPosition(s string) (int, error) {
	letter := []rune(s)[0]
	first := []rune("a")[0]
	if letter < first {
		return 0, errors.New("column (lettered index) out of range!")
	}
	return int(letter - first), nil
}

func letterInPosition(n int) string {
	first := []rune("a")[0]
	return string(first + rune(n))
}

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

	dim := dimensions{10, 10}

	genericShips := canonicalBattleship()

	attackField := makeScatteredField(dim, genericShips, adversary)
	defendField := makeScatteredField(dim, genericShips, human)

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
