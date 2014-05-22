package main

import (
    "fmt"
    "strconv"
    "math/rand"
    "time"
    "strings"
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
        var dummyField field
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
            // set up a dummy field full of our own ships
            dummyField = field{d, ships, []coord{}}
            for i := 0; i < s.length; i++ {
                if horizontal {
                    coords[i] = c.right(i)
                    occupied, _ = dummyField.shoot(c.right(i))
                } else {
                    coords[i] = c.down(i)
                    occupied, _ = dummyField.shoot(c.down(i))
                }
                // if a dummy shot hits, the space is occupied
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

func move(f field) (string, field) {
    var userInput string
    fmt.Scanln("%s", &userInput)
    rawCoord := []rune(userInput)
    rowLetter := strings.ToLower(string(rawCoord[0:1]))
    column, err := strconv.Atoi(string(rawCoord[1:]))
    if err != nil {
        panic(err)
    }
    hit, hitShip := f.shoot(coord{row, column})
    if hit {
        // FIXME: pointers!!!!!!!!!!! TODO
    } else {
        return "Miss!", f
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    dim := dimensions{10, 10}

    genericShips = canonicalBattleship()

    adversaryShips := makeShips(dim, genericShips, adversary)
    // TODO: Let humans choose
    humanShips := makeShips(dim, genericShips, human)

    field := field{
        dimensions: dim,
        ships: append(adversaryShips, humanShips),
        misses: []coord{},
    }

    winner := field.winner()
    for winner == nobody {
        fmt.Print(field.humanView())
        // game loop
    }
    if winner == human {
        fmt.Println("You've won! Congratulations.")
    }
}
