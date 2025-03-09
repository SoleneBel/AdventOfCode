# Advent of Code 2022 - Day 17: Pyroclastic Flow

## Description
This project is a solution to the Day 17 problem of Advent of Code 2022, titled "Pyroclastic Flow". 
It simulates the fall of rocks of various shapes into a fixed-width chamber, considering lateral movements and gravity. 
The goal is to determine the final height of the rock pile after an extremely large number of iterations.

## Features
- Loads movement instructions from an input file.
- Simulates the falling and lateral movement of rocks.
- Detects and exploits cycles to speed up execution.

## Installation
1. Ensure you have Go installed on your machine.
2. Clone this repository.
3. Compile and run the program:.

## Usage
The program expects an input file containing the rock movement instructions.
By default, it uses the file `inputTest.txt`, which is the example given.

```sh
go run main.go
```

If you want to use a different input file, modify the call to `pyroclasticFlow("input.txt")` in `main.go`.

## Code Explanation
The program is structured as follows:
- **Main Structures**
    - `State`: Represents a simulation state for cycle detection.
    - `Chamber`: Represents the chamber where rocks fall.
- **Key Functions**
    - `addRock(c *Chamber, r []string)`: Adds a rock to the chamber.
    - `rockFall(c *Chamber, r []string, f string, ind int) int`: Simulates the fall of a rock.
    - `checkCycle(...)`: Detects cycles to accelerate the simulation.
    - `pyroclasticFlow(f string)`: Main function running the simulation.

## Author
Developed by Solène Bellissard.

## Resources
- [Advent of Code 2022 - Day 17](https://adventofcode.com/2022/day/17)

