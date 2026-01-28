package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Shape struct {
	cells  [][]bool
	width  int
	height int
}

type Region struct {
	width    int
	height   int
	presents []int
}

func parseShape(lines []string) Shape {
	if len(lines) == 0 {
		return Shape{}
	}

	height := len(lines)
	width := len(lines[0])

	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
		for j := 0; j < width && j < len(lines[i]); j++ {
			cells[i][j] = (lines[i][j] == '#')
		}
	}

	return Shape{cells: cells, width: width, height: height}
}

func rotate90(shape Shape) Shape {
	newHeight := shape.width
	newWidth := shape.height
	newCells := make([][]bool, newHeight)

	for i := 0; i < newHeight; i++ {
		newCells[i] = make([]bool, newWidth)
		for j := 0; j < newWidth; j++ {
			newCells[i][j] = shape.cells[newWidth-1-j][i]
		}
	}

	return Shape{cells: newCells, width: newWidth, height: newHeight}
}

func flipHorizontal(shape Shape) Shape {
	newCells := make([][]bool, shape.height)
	for i := 0; i < shape.height; i++ {
		newCells[i] = make([]bool, shape.width)
		for j := 0; j < shape.width; j++ {
			newCells[i][j] = shape.cells[i][shape.width-1-j]
		}
	}
	return Shape{cells: newCells, width: shape.width, height: shape.height}
}

func shapeKey(shape Shape) string {
	var sb strings.Builder
	for i := 0; i < shape.height; i++ {
		for j := 0; j < shape.width; j++ {
			if shape.cells[i][j] {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func getTransformations(shape Shape) []Shape {
	seen := make(map[string]bool)
	var unique []Shape

	current := shape
	for i := 0; i < 4; i++ {
		key := shapeKey(current)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, current)
		}
		current = rotate90(current)
	}

	current = flipHorizontal(shape)
	for i := 0; i < 4; i++ {
		key := shapeKey(current)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, current)
		}
		current = rotate90(current)
	}

	return unique
}

func canPlace(grid [][]bool, shape Shape, row, col int) bool {
	if row+shape.height > len(grid) || col+shape.width > len(grid[0]) {
		return false
	}

	for i := 0; i < shape.height; i++ {
		for j := 0; j < shape.width; j++ {
			if shape.cells[i][j] && grid[row+i][col+j] {
				return false
			}
		}
	}
	return true
}

func place(grid [][]bool, shape Shape, row, col int) {
	for i := 0; i < shape.height; i++ {
		for j := 0; j < shape.width; j++ {
			if shape.cells[i][j] {
				grid[row+i][col+j] = true
			}
		}
	}
}

func remove(grid [][]bool, shape Shape, row, col int) {
	for i := 0; i < shape.height; i++ {
		for j := 0; j < shape.width; j++ {
			if shape.cells[i][j] {
				grid[row+i][col+j] = false
			}
		}
	}
}

func countCells(shape Shape) int {
	count := 0
	for i := 0; i < shape.height; i++ {
		for j := 0; j < shape.width; j++ {
			if shape.cells[i][j] {
				count++
			}
		}
	}
	return count
}

func tryFit(grid [][]bool, shapes []Shape, presentsToPlace []int, presentIdx int, transforms [][]Shape) bool {
	if presentIdx == len(presentsToPlace) {
		return true
	}

	shapeIdx := presentsToPlace[presentIdx]
	transformations := transforms[shapeIdx]

	for _, transform := range transformations {
		for row := 0; row <= len(grid)-transform.height; row++ {
			for col := 0; col <= len(grid[0])-transform.width; col++ {
				if canPlace(grid, transform, row, col) {
					place(grid, transform, row, col)
					if tryFit(grid, shapes, presentsToPlace, presentIdx+1, transforms) {
						return true
					}
					remove(grid, transform, row, col)
				}
			}
		}
	}

	return false
}

func canFitPresents(region Region, shapes []Shape, transforms [][]Shape) bool {
	grid := make([][]bool, region.height)
	for i := range grid {
		grid[i] = make([]bool, region.width)
	}

	var presentsToPlace []int
	totalCells := 0
	for shapeIdx, count := range region.presents {
		if shapeIdx >= len(shapes) {
			continue
		}
		for i := 0; i < count; i++ {
			presentsToPlace = append(presentsToPlace, shapeIdx)
			totalCells += countCells(shapes[shapeIdx])
		}
	}

	// Quick check: if total cells exceed grid size, impossible
	if totalCells > region.width*region.height {
		return false
	}

	return tryFit(grid, shapes, presentsToPlace, 0, transforms)
}

func main() {
	file, err := os.Open("inputday12.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var shapes []Shape
	var regions []Region

	// Parse file
	var currentPattern []string
	inShapeSection := true

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line is a region line (format: "NNxNN: ...")
		if strings.Contains(line, "x") && strings.Contains(line, ":") {
			colonIdx := strings.Index(line, ":")
			beforeColon := line[:colonIdx]
			if strings.Contains(beforeColon, "x") {
				// This is a region line
				if len(currentPattern) > 0 {
					shapes = append(shapes, parseShape(currentPattern))
					currentPattern = nil
				}
				inShapeSection = false

				// Parse the region
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					dims := strings.Split(parts[0], "x")
					if len(dims) == 2 {
						width, _ := strconv.Atoi(dims[0])
						heightStr := strings.TrimSuffix(dims[1], ":")
						height, _ := strconv.Atoi(heightStr)

						var presents []int
						for i := 1; i < len(parts); i++ {
							count, _ := strconv.Atoi(parts[i])
							presents = append(presents, count)
						}

						regions = append(regions, Region{width: width, height: height, presents: presents})
					}
				}
				continue
			}
		}

		if inShapeSection {
			if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
				// Shape header - save previous pattern
				if len(currentPattern) > 0 {
					shapes = append(shapes, parseShape(currentPattern))
					currentPattern = nil
				}
			} else if line != "" {
				// Part of shape pattern
				currentPattern = append(currentPattern, line)
			} else if len(currentPattern) > 0 {
				// Empty line after pattern
				shapes = append(shapes, parseShape(currentPattern))
				currentPattern = nil
			}
		}
	}

	// Save last pattern if exists
	if len(currentPattern) > 0 {
		shapes = append(shapes, parseShape(currentPattern))
	}

	fmt.Printf("Parsed %d shapes and %d regions\n", len(shapes), len(regions))

	// Precompute all transformations
	transforms := make([][]Shape, len(shapes))
	for i := range shapes {
		transforms[i] = getTransformations(shapes[i])
	}

	count := 0
	for i, region := range regions {
		if canFitPresents(region, shapes, transforms) {
			count++
		}
		if (i+1)%100 == 0 || i < 10 {
			fmt.Printf("Progress: %d/%d regions checked, %d fit\n", i+1, len(regions), count)
		}
	}

	fmt.Printf("\nAnswer: %d\n", count)
}
