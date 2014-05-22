package main

func canonicalBattleship() [5]ship {
    return [...]ship{
        ship{
            name: "Aircraft Carrier",
            length: 5,
            spaces: make([]coord, 5),
            holes: make([]coord, 5),
            owner: nobody,
        },
        ship{
            name: "Battleship",
            length: 4,
            spaces: make([]coord, 4),
            holes: make([]coord, 4),
            owner: nobody,
        },
        ship{
            name: "Submarine",
            length: 3,
            spaces: make([]coord, 3),
            holes: make([]coord, 3),
            owner: nobody,
        },
        ship{
            name: "Cruiser",
            length: 3,
            spaces: make([]coord, 3),
            holes: make([]coord, 3),
            owner: nobody,
        },
        ship{
            name: "Patrol Boat",
            length: 2,
            spaces: make([]coord, 2),
            holes: make([]coord, 2),
            owner: nobody,
        },
    }
}
