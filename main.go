package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const WIDTH = 7
const ITERATIONS = 1000000000001 // 2023 or 1000000000001

// Rocks' shapes
var shapes = [][]string{
	{
		"@@",
		"@@",
	},
	{
		"@@@@",
	},
	{
		".@.",
		"@@@",
		".@.",
	},
	{
		"@@@",
		"..@",
		"..@",
	},
	{
		"@",
		"@",
		"@",
		"@",
	},
}

// State : Struct for cycles' states
type State struct {
	loop   int // L'indice de boucle
	height int // La taille de la chambre
}

// Chamber : the chamber where the rocks fall
type Chamber struct {
	width  int
	height int
	grid   [][]byte
}

// Get the input
func parseFile(filename string) string {
	var f, _ = os.Open(filename)
	sc := bufio.NewScanner(f)

	// Read the line of the input
	for sc.Scan() {
		line := sc.Text()

		if line != "" {
			return line
		}
	}
	return ""
}

// Initialize the chamber's grid with '.' (means there's nothing)
func initGrid(width, height int) [][]byte {
	grid := make([][]byte, height)
	for i := range grid {
		grid[i] = make([]byte, width)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	return grid
}

// Initialize the chamber
func initChamber() Chamber {
	return Chamber{
		width:  WIDTH,
		height: 3,
		grid:   initGrid(WIDTH, 3),
	}
}

// Add and center a rock in the chamber
func addRock(c *Chamber, r []string) {
	grid := c.grid
	width := c.width

	for _, row := range r {
		newLine := make([]byte, width)

		// Fill the line with '.'
		for i := range newLine {
			newLine[i] = '.'
		}

		// Place the left side of the rock on the 3rd column
		for i, char := range row {
			newLine[2+i] = byte(char)
		}

		// Add the line to the grid
		grid = append(grid, newLine)
	}

	// Chamber's update
	c.grid = grid
	c.height += len(r)
}

// Check if a rock can fall
func canFall(grid [][]byte) bool {
	for i := 0; i < len(grid)-1; i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == '@' && (i == 0 || grid[i-1][j] == '#') {
				return false // If there's a collision with a wall or a rock (#)
			}
		}
	}
	return true
}

// Check if a rock can go left or right
func canShift(grid [][]byte, direction uint8) bool {
	for i := range grid {
		length := len(grid[i])
		for j := range grid[i] {
			// Below is a horrible line to say that if we are against a wall or next to a #, we cannot shift the rock.
			if grid[i][j] == '@' && (direction == '>' && (j == length-1 || grid[i][j+1] == '#') || direction == '<' && (j == 0 || grid[i][j-1] == '#')) {
				return false
			}
		}
	}
	return true
}

// Check if a line contains a part of a rock
// typeRock can be chosen between '#' (rock at rest) pr '@' (rock falling)
func containsRock(line []byte, typeRock byte) bool {
	for _, char := range line {
		if char == typeRock {
			return true
		}
	}
	return false
}

// Move right or left the rock regarding the movement to follow
func shiftRock(grid [][]byte, direction uint8) {
	for i := range grid {
		if containsRock(grid[i], '@') { // We shift only if we have a rock '@'
			if direction == '>' { // If right
				for j := len(grid[i]) - 1; j > 0; j-- {
					if grid[i][j] == '.' && grid[i][j-1] == '@' {
						grid[i][j] = '@'
						grid[i][j-1] = '.'
					}
				}
			} else if direction == '<' { // If left
				for j := 0; j < len(grid[i])-1; j++ {
					if grid[i][j] == '.' && grid[i][j+1] == '@' {
						grid[i][j] = '@'
						grid[i][j+1] = '.'
					}
				}
			} // If there are other symbols (meaning wrong input), they're ignored
		}
	}
}

// Simulate the fall of a rock from one step
func fallOneStep(grid [][]byte) {
	for i := range grid {
		if containsRock(grid[i], '@') {
			for j := range grid[i] {
				if grid[i][j] == '@' {
					grid[i-1][j] = '@'
					grid[i][j] = '.'
				}
			}
		}
	}
}

// Once the rock can't fall more, it "rests" (going from "@" to "#")
// Useful for the collision tests
func restRock(grid [][]byte) {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == '@' {
				grid[i][j] = '#'
			}
		}
	}
}

// Check if a line is empty
func emptyLine(line []byte) bool {
	for _, char := range line {
		if char != '.' {
			return false
		}
	}
	return true
}

// Count the number of empty lines at the top of the chamber (the end of the grid)
func countEmptyLinesAtEnd(c Chamber) int {
	counter := 0

	for i := len(c.grid) - 1; i >= 0; i-- {
		if emptyLine(c.grid[i]) {
			counter++
		} else {
			break
		}
	}

	return counter
}

// Put exactly three empty lines at the top of the chamber if needed
// Useful bc some rocks weren't added at the right height
func fixEmptyLinesAtEnd(c *Chamber) {
	emptyCount := countEmptyLinesAtEnd(*c)

	if emptyCount > 3 { // Delete extra lines
		c.grid = c.grid[:len(c.grid)-(emptyCount-3)]
	} else if emptyCount < 3 { // Add necessary lines
		width := c.width
		emptyRow := make([]byte, width)

		for i := range emptyRow {
			emptyRow[i] = '.'
		}

		for i := 0; i < 3-emptyCount; i++ {
			c.grid = append(c.grid, emptyRow)
		}
	}

	c.height = len(c.grid)
}

// Simulate the fall of a rock in the chamber
func rockFall(c *Chamber, r []string, f string, ind int) int {
	addRock(c, r)
	grid := c.grid
	i := c.height - 1 // We begin from the top of the chamber (= the end of the grid)
	index := ind      // We resume the movements from where we stopped
	nbMovements := len(f)

	for i >= 0 {
		movement := f[index%nbMovements] // get '<' or '>'
		index++

		if canShift(grid, movement) {
			shiftRock(grid, movement)
		}

		if !canFall(grid) {
			break
		}

		fallOneStep(grid)

		i--
	}

	restRock(grid)

	c.grid = grid // Grid's update

	return index
}

// Remove the lines only with ..... to have only the interesting part of the chamber
func removeTopChamber(c *Chamber) {
	var newGrid [][]byte

	// Keeping only the lines with at least a "#" or a "@"
	for _, row := range c.grid {
		if containsRock(row, '#') || containsRock(row, '@') {
			newGrid = append(newGrid, row)
		}
	}

	c.grid = newGrid
	c.height = len(newGrid)
}

////////////////
// FOR PART 2
////////////////

// Check if every column has at least one '#'
func fullCoverageWidth(grid [][]byte) bool {
	found := false

	for i := range grid {
		for _, symb := range grid[i] {
			if symb == '#' {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// Get top of the chamber
func getTopChamber(c Chamber) [][]byte {
	i := len(c.grid) - 1

	for i > 0 {
		if !emptyLine(c.grid[i]) {
			break
		}
		i--
	}

	if i < 8 {
		return nil
	}

	topChamber := c.grid[i-7 : i]

	if fullCoverageWidth(topChamber) {
		return topChamber
	}

	return nil
}

// Create a key to stock our state
func encodeKey(grid [][]byte, iShape, iMovement int) string {
	var sb strings.Builder

	// Encoding using string
	for _, row := range grid {
		sb.Write(row)
	}

	// Adding the indexes
	sb.WriteByte('|')
	sb.WriteString(strconv.Itoa(iShape))
	sb.WriteByte(',')
	sb.WriteString(strconv.Itoa(iMovement))

	return sb.String()
}

// Check if there's an existing cycle
func checkCycle(c *Chamber, stateMap *map[string]State, iShape, iLoop, iMovement int, foundCycle bool) (int, int, bool) {
	topChamber := getTopChamber(*c)
	difHeight := 0

	if topChamber != nil {
		key := encodeKey(topChamber, iShape, iMovement) // Encoding key

		previous := (*stateMap)[key] // Previous state with this same encoding

		// Update or Add new state
		state := (*stateMap)[key]
		state.loop = iLoop
		state.height = c.height
		(*stateMap)[key] = state

		// Calculate cycle and height
		if previous.loop != 0 && previous.height != 0 { // If a previous state existed, then there's a cycle
			cycle := iLoop - previous.loop
			// fmt.Println("cycle:", iLoop, previous.loop, cycle)
			nbCyclesTot := (ITERATIONS-iLoop)/cycle - 1 // Total number of cycles

			difHeight = (c.height - previous.height) * nbCyclesTot // Gets the height
			iLoop += nbCyclesTot * cycle                           // Skipping t othe end of the last cycle
			foundCycle = true
		}
	}

	return difHeight, iLoop, foundCycle
}

////////////////
// MAIN FUNCTION
////////////////

// Main function
func pyroclasticFlow(f string) {
	file := parseFile(f)

	// Init of structures and variables
	c := initChamber()
	stateMap := make(map[string]State)
	foundCycle := false
	nbRock := len(shapes)
	nbMovements := len(file)
	difHeight := 0 // Height gained with cycles
	ind := 0       // Index in the list of movements
	i := 1         // Iterations

	for i < ITERATIONS { // 2022 for pat 1 or 1000000000001 for part 2
		fixEmptyLinesAtEnd(&c)

		// Get the index of the rock that will fall
		index := i % nbRock

		ind = rockFall(&c, shapes[index], file, ind)

		if i >= nbMovements && !foundCycle {
			difHeight, i, foundCycle = checkCycle(&c, &stateMap, index, i, ind%nbMovements, foundCycle)
		}

		i++
	}

	// All the empty lines at the top of the chamber are removed to have the correct height
	removeTopChamber(&c)

	fmt.Println("Chamber's height : ", c.height+difHeight, "!")

}

func main() {
	pyroclasticFlow("inputTest.txt")
}
