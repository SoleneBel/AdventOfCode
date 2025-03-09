package main

import (
	"testing"
)

func TestInitChamber(t *testing.T) {
	chamber := initChamber()
	if chamber.width != WIDTH {
		t.Errorf("Expected chamber width to be %d, got %d", WIDTH, chamber.width)
	}
	if chamber.height != 3 {
		t.Errorf("Expected chamber height to be 3, got %d", chamber.height)
	}
}

func TestAddRock(t *testing.T) {
	chamber := initChamber()
	rock := shapes[0] // A simple square rock
	addRock(&chamber, rock)

	if chamber.height <= 3 {
		t.Errorf("Expected chamber height to increase after adding a rock")
	}

	// Check if the rock is correctly placed starting from the 3rd column
	for i, row := range rock {
		for j, char := range row {
			if chamber.grid[len(chamber.grid)-len(rock)+i][j+2] != byte(char) {
				t.Errorf("Rock was not placed correctly in the chamber at the expected position")
			}
		}
	}
}

func TestCanFall(t *testing.T) {
	chamber := initChamber()
	rock := shapes[0]
	addRock(&chamber, rock)
	if !canFall(chamber.grid) {
		t.Errorf("Rock should be able to fall initially")
	}
}

func TestCanShift(t *testing.T) {
	chamber := initChamber()
	rock := shapes[0]
	addRock(&chamber, rock)

	if !canShift(chamber.grid, '<') {
		t.Errorf("Rock should be able to shift left initially")
	}
	if !canShift(chamber.grid, '>') {
		t.Errorf("Rock should be able to shift right initially")
	}
}

func TestRockFall(t *testing.T) {
	chamber := initChamber()
	rock := shapes[0]
	movements := "<>>><<"
	index := rockFall(&chamber, rock, movements, 0)

	if index == 0 {
		t.Errorf("Rock fall did not progress movement index")
	}
}

func TestCycleDetection(t *testing.T) {
	chamber := initChamber()
	stateMap := make(map[string]State)
	foundCycle := false
	difHeight, iLoop, foundCycle := checkCycle(&chamber, &stateMap, 0, 10, 5, foundCycle)

	if foundCycle && difHeight == 0 {
		t.Errorf("Cycle detection should modify height difference")
	}
	if foundCycle && iLoop == 10 {
		t.Errorf("Cycle detection should alter loop count")
	}
}
