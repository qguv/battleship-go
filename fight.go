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

func randomCoord(d dimensions) (c coord) {
	c.x = rand.Intn(d.x)
	c.y = rand.Intn(d.y)
	return
}

func makeShips(d dimensions, generics []ship, owner player) []ship {
	ships := make([]ship, len(generics))
	for shipI, s := range generics {
		var horizontal bool
		var c coord
		var lastc coord
		var coords []coord
		var occupied bool
	TryCoords:
		for {
			c = randomCoord(d)
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
		s.owner = adversary
		ships[shipI] = s
	}
	return ships
}

func alphabetPosition(s string) (int, error) {
	letter := []rune(s)[0]
	first := []rune("a")[0]
	if letter < first {
		return 0, errors.New("column (lettered index) out of range!")
	}
	return int(letter - first), nil
}

func move(f *field) string {
	var userInput string
	fmt.Scanln("%s", &userInput)
	rawCoord := []rune(userInput)
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
		return fmt.Sprintln("you hit a", hitShip.name)
	} else {
		return "Miss!"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dim := dimensions{10, 10}

	genericShips := canonicalBattleship()

	adversaryShips := makeShips(dim, genericShips, adversary)
	// TODO: Let humans choose
	humanShips := makeShips(dim, genericShips, human)

	field := field{
		dimensions: dim,
		ships:      append(adversaryShips, humanShips...),
		misses:     []coord{},
	}

	winner := field.winner()
	for winner == nobody {
		// game loop
	}
	if winner == human {
		fmt.Println("You've won! Congratulations.")
	}
}
