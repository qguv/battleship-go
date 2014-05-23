package main

import "testing"

func TestCoordDimensions(t *testing.T) {
	origin := coord{0, 0}

	c0 := coord{3, 5}
	obs0 := origin.dimensions(c0)
	exp0 := dimensions{3, 5}
	if obs0 != exp0 {
		t.Errorf("dimensions method on coord (against origin) returned incorrect result x=%v, y=%v", obs0.x, obs0.y)
	}

	c1 := coord{7, 6}
	obs1 := c0.dimensions(c1)
	exp1 := dimensions{4, 1}
	if obs1 != exp1 {
		t.Error("dimensions method on coord (against another coord) returned incorrect result")
	}

	c2 := coord{7, 3}
	obs2 := c0.dimensions(c2)
	exp2 := dimensions{4, -2}
	if obs2 != exp2 {
		t.Error("negative-producing dimensions method on coord returned incorrect result")
	}
}

func TestCoordTransforms(t *testing.T) {
	origin := coord{0, 0}
	expDown := coord{0, 2}
	if origin.down(2) != expDown {
		t.Error("down method on coord returned incorrect result")
	}

	expRight := coord{2, 0}
	if origin.right(2) != expRight {
		t.Error("right method on coord returned incorrect result")
	}
}

func TestArea(t *testing.T) {
	c0 := coord{0, 0}
	c1 := coord{5, 5}
	dim := dimensions{5, 5}

	if c0.area(c1) != 25 {
		t.Error("area method on coord returned incorrect result")
	}

	if dim.area() != 25 {
		t.Error("area method on dimensions returned incorrect result")
	}
}

func TestWithin(t *testing.T) {
	bounds := dimensions{5, 5}
	goods := []coord{
		coord{4, 0},
		coord{2, 3},
		coord{0, 0},
		coord{1, 4},
	}
	bads := []coord{
		coord{0, -1},
		coord{-1, 0},
		coord{-1, -1},
		coord{5, 5},
		coord{5, 4},
		coord{4, 5},
		coord{20, 3},
		coord{3, 20},
		coord{20, 20},
	}

	for _, c := range goods {
		if !c.within(bounds) {
			t.Fatal("a coordinate in-bounds was reported out-of-bounds by its within method")
		}
	}
	for _, c := range bads {
		if c.within(bounds) {
			t.Fatal("a coordinate out-of-bounds was reported in-bounds by its within method")
		}
	}
}

func TestShipIsDestroyed(t *testing.T) {
	prestine := ship{
		name:   "A Good-looking Ship",
		length: 4,
		spaces: []coord{
			coord{2, 3},
			coord{2, 4},
			coord{2, 5},
			coord{2, 6},
		},
		holes: make([]coord, 4),
		owner: human,
	}
	damaged := ship{
		name:   "An Experienced Ship",
		length: 4,
		spaces: []coord{
			coord{2, 3},
			coord{2, 4},
			coord{2, 5},
			coord{2, 6},
		},
		holes: []coord{
			coord{2, 4},
			coord{2, 5},
			coord{2, 6},
		},
		owner: human,
	}
	destroyed := ship{
		name:   "A Wasted Ship",
		length: 4,
		spaces: []coord{
			coord{2, 3},
			coord{2, 4},
			coord{2, 5},
			coord{2, 6},
		},
		holes: []coord{
			coord{2, 3},
			coord{2, 4},
			coord{2, 5},
			coord{2, 6},
		},
		owner: human,
	}
	if prestine.isDestroyed() {
		t.Error("a prestine ship's isDestroyed method erroneously reports its demise")
	}
	if damaged.isDestroyed() {
		t.Error("a damaged ship's isDestroyed method erroneously reports its demise")
	}
	if !destroyed.isDestroyed() {
		t.Error("a destroyed ship's isDestroyed method fails to reports its demise")
	}
}

func TestShipsLeft(t *testing.T) {
	allDestroyed := field{
		dimensions: dimensions{2, 2},
		misses:     []coord{},
		ships: []ship{
			ship{
				name:   "Baddie",
				length: 1,
				spaces: []coord{
					coord{1, 1},
				},
				holes: []coord{
					coord{1, 1},
				},
				owner: adversary,
			},
		},
	}
	if allDestroyed.shipsLeft() {
		t.Error("shipsLeft reporting remaining ships in a destroyed field")
	}

	stillKicking := field{
		dimensions: dimensions{2, 2},
		misses:     []coord{},
		ships: []ship{
			ship{
				name:   "Baddie",
				length: 1,
				spaces: []coord{
					coord{1, 1},
				},
				holes: []coord{},
				owner: adversary,
			},
		},
	}
	if !stillKicking.shipsLeft() {
		t.Error("shipsLeft reporting no remaining ships in an active field")
	}
}

func TestShoot(t *testing.T) {
	dims := dimensions{3, 3}
	ships := []ship{
		ship{
			name:   "Enemy Sailboat",
			length: 1,
			spaces: []coord{
				coord{1, 1},
			},
			holes: []coord{},
			owner: adversary,
		},
		ship{
			name:   "Our Sailboat",
			length: 1,
			spaces: []coord{
				coord{1, 2},
			},
			holes: []coord{},
			owner: human,
		},
	}
	var misses []coord
	f := field{
		dimensions: dims,
		misses:     misses,
		ships:      ships,
	}

	badHit, _ := f.shoot(coord{0, 1})
	if badHit {
		t.Error("shooting an unoccupied coordinate in a field resulted in a hit")
	}

	goodHit, ship := f.shoot(coord{1, 1})
	if !goodHit {
		t.Error("shooting an occupied coordinate in a field resulted in a miss")
	}
	if ship.name != f.ships[0].name {
		t.Error("shooting an occupied coordinate in a field resulted in the wrong ship being returned")
	}
	if !ship.isDestroyed() {
		t.Error("shooting a ship completely did not result in its destruction")
	}
}

func TestOccupied(t *testing.T) {
	ships := []ship{
		ship{
			name:   "Tyrone",
			length: 1,
			spaces: []coord{
				coord{0, 0},
			},
			holes: []coord{},
			owner: human,
		},
	}

	goodOccupied := coordOccupied(coord{0, 0}, ships)
	if !goodOccupied {
		t.Error("coordOccupied returns a false negative")
	}

	badOccupied := coordOccupied(coord{0, 1}, ships)
	if badOccupied {
		t.Error("coordOccupied returns a false positive")
	}
}
