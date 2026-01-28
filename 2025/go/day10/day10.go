package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("inputday10.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	configPresses := 0
	joltPresses := 0.0

	for scanner.Scan() {
		line := scanner.Text()
		machine, err := deserialize(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deserializing line: %v\n", err)
			continue
		}

		presses, err := machine.configure()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error configuring machine: %v\n", err)
			continue
		}
		configPresses += presses

		presses2, err := machine.jolt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error jolting machine: %v\n", err)
			continue
		}
		joltPresses += presses2
	}

	fmt.Printf("Part 1: %d\n", configPresses)
	fmt.Printf("Part 2: %.0f\n", joltPresses)
}

type Machine struct {
	configuration int
	buttons       []int
	joltage       []int
}

func (m *Machine) createMatrix() [][]float64 {
	matrix := make([][]float64, len(m.joltage))
	ptr := 1
	for i := range matrix {
		matrix[i] = make([]float64, len(m.buttons)+1)
		matrix[i][len(m.buttons)] = float64(m.joltage[i])
		for j, button := range m.buttons {
			if button&ptr > 0 {
				matrix[i][j] = 1
			} else {
				matrix[i][j] = 0
			}
		}
		ptr = ptr << 1
	}
	return matrix
}

func getMaxPresses(m *Machine, matrix [][]float64) map[int]int {
	maxPresses := make(map[int]int)

	for k := 0; k < len(matrix[0])-1; k++ {
		for h := 0; h < len(matrix); h++ {
			if matrix[h][k] == 1 {
				if _, ok := maxPresses[m.buttons[k]]; !ok {
					maxPresses[m.buttons[k]] = m.joltage[h]
				} else {
					if maxPresses[m.buttons[k]] > m.joltage[h] {
						maxPresses[m.buttons[k]] = m.joltage[h]
					}
				}
			}
		}
	}
	return maxPresses
}

func (m *Machine) jolt() (float64, error) {
	matrix := m.createMatrix()
	maxPresses := getMaxPresses(m, matrix)

	rref := matrixReduce(matrix)
	freeVariables, params := parametrize(rref)

	ranges := make([]rangeStruct, len(freeVariables)+1)
	ranges[0] = rangeStruct{min: 1, max: 2}
	for i, j := range freeVariables {
		if _, ok := maxPresses[m.buttons[i]]; !ok {
			ranges[j+1] = rangeStruct{min: 0, max: 1}
		} else {
			ranges[j+1] = rangeStruct{min: 0, max: maxPresses[m.buttons[i]] + 1}
		}
	}

	combs := generateCombinations(ranges)

	var best *float64
	for _, comb := range combs {
		floatComb := make([]float64, len(comb))
		for i, v := range comb {
			floatComb[i] = float64(v)
		}

		coefs := getCoefficients(params, floatComb)
		if roundWhileValidatingCoefs(coefs) && coefsConsistentWithMatrix(matrix, coefs) {
			acc := 0.0
			for _, coef := range coefs {
				acc += coef
			}
			if best == nil || acc < *best {
				best = &acc
			}
		}
	}

	if best == nil {
		return 0, fmt.Errorf("no combinations found")
	}

	return *best, nil
}

func (m *Machine) configure() (int, error) {
	combs := calcCombs(len(m.buttons))
	var best *int

	for i := 0; i <= combs; i++ {
		val, presses := getValPresses(i, m)
		if val == m.configuration {
			if best == nil || *best > presses {
				best = &presses
			}
		}
	}

	if best == nil {
		return 0, fmt.Errorf("no configuration found")
	}

	return *best, nil
}

func getValPresses(comb int, m *Machine) (int, int) {
	val := 0
	presses := 0
	ptr := 1
	button := 0
	for ptr <= comb {
		if comb&ptr > 0 {
			val = press(val, m.buttons[button])
			presses++
		}
		ptr = ptr << 1
		button++
	}
	return val, presses
}

func calcCombs(buttons int) int {
	combs := 1
	for range buttons {
		combs = combs << 1
	}
	combs = combs - 1
	return combs
}

func deserialize(s string) (*Machine, error) {
	var builder strings.Builder
	configuration := 0
	ptr := 1
	buttons := make([]int, 0)
	var button int
	joltage := make([]int, 0)

	state := 0 // 0=start, 1=brackets, 2=after brackets, 3=parens, 4=after parens, 5=braces

	for _, r := range s {
		switch r {
		case '[':
			state = 1
		case ']':
			state = 2
		case '(':
			state = 3
		case ')':
			if builder.Len() > 0 {
				res, _ := strconv.Atoi(builder.String())
				builder.Reset()
				button = addIndicator(button, res)
			}
			buttons = append(buttons, button)
			button = 0
			state = 4
		case '{':
			state = 5
		case '}':
			if builder.Len() > 0 {
				res, _ := strconv.Atoi(builder.String())
				builder.Reset()
				joltage = append(joltage, res)
			}
		case ' ':
			continue
		case ',':
			if state == 3 {
				res, _ := strconv.Atoi(builder.String())
				builder.Reset()
				button = addIndicator(button, res)
			} else if state == 5 {
				res, _ := strconv.Atoi(builder.String())
				builder.Reset()
				joltage = append(joltage, res)
			}
		case '.':
			if state == 1 {
				ptr = ptr << 1
			}
		case '#':
			if state == 1 {
				configuration = ptr | configuration
				ptr = ptr << 1
			}
		default:
			if state == 3 || state == 5 {
				builder.WriteRune(r)
			}
		}
	}

	return &Machine{
		configuration: configuration,
		buttons:       buttons,
		joltage:       joltage,
	}, nil
}

func addIndicator(button, indicator int) int {
	x := 1
	for range indicator {
		x = x << 1
	}
	return button | x
}

func press(curr, button int) int {
	return (curr | button) ^ (curr & button)
}

// Math utilities

type rangeStruct struct {
	min, max int
}

func isZero(f float64) bool {
	return math.Abs(f) < 1e-9
}

func matrixReduce(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])

	// Create copy
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		copy(m[i], matrix[i])
	}

	// Gaussian elimination with partial pivoting
	lead := 0
	for r := 0; r < rows; r++ {
		if lead >= cols {
			break
		}

		// Find pivot
		i := r
		for isZero(m[i][lead]) {
			i++
			if i == rows {
				i = r
				lead++
				if lead == cols {
					return m
				}
			}
		}

		// Swap rows
		m[i], m[r] = m[r], m[i]

		// Normalize pivot row
		if !isZero(m[r][lead]) {
			div := m[r][lead]
			for j := 0; j < cols; j++ {
				m[r][j] /= div
			}
		}

		// Eliminate column
		for i := 0; i < rows; i++ {
			if i != r {
				mult := m[i][lead]
				for j := 0; j < cols; j++ {
					m[i][j] -= mult * m[r][j]
				}
			}
		}
		lead++
	}

	return m
}

func parametrize(rref [][]float64) (map[int]int, [][]float64) {
	rows := len(rref)
	cols := len(rref[0]) - 1 // Exclude augmented column

	pivotCols := make([]int, 0)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if !isZero(rref[r][c]) && math.Abs(rref[r][c]-1.0) < 1e-9 {
				pivotCols = append(pivotCols, c)
				break
			}
		}
	}

	freeVars := make(map[int]int)
	freeIdx := 0
	for c := 0; c < cols; c++ {
		isPivot := false
		for _, pc := range pivotCols {
			if pc == c {
				isPivot = true
				break
			}
		}
		if !isPivot {
			freeVars[c] = freeIdx
			freeIdx++
		}
	}

	// Build parametric solution
	params := make([][]float64, cols)
	for i := range params {
		params[i] = make([]float64, len(freeVars)+1)
	}

	for c := 0; c < cols; c++ {
		if _, isFree := freeVars[c]; isFree {
			params[c][freeVars[c]+1] = 1.0
		} else {
			// Find row with pivot in this column
			for r := 0; r < rows; r++ {
				if !isZero(rref[r][c]) && math.Abs(rref[r][c]-1.0) < 1e-9 {
					params[c][0] = rref[r][cols]
					for fc, fi := range freeVars {
						params[c][fi+1] = -rref[r][fc]
					}
					break
				}
			}
		}
	}

	return freeVars, params
}

func getCoefficients(params [][]float64, values []float64) []float64 {
	coefs := make([]float64, len(params))
	for i := range params {
		val := params[i][0]
		for j := 1; j < len(params[i]); j++ {
			val += params[i][j] * values[j]
		}
		coefs[i] = val
	}
	return coefs
}

func roundWhileValidatingCoefs(coefs []float64) bool {
	for i, coef := range coefs {
		rounded := math.Round(coef)
		if !isZero(coef-rounded) || rounded < 0 {
			return false
		}
		coefs[i] = rounded
	}
	return true
}

func coefsConsistentWithMatrix(matrix [][]float64, coefs []float64) bool {
	for i := 0; i < len(matrix); i++ {
		sum := 0.0
		for j := 0; j < len(coefs); j++ {
			sum += matrix[i][j] * coefs[j]
		}
		expected := matrix[i][len(matrix[i])-1]
		if !isZero(sum - expected) {
			return false
		}
	}
	return true
}

func generateCombinations(ranges []rangeStruct) [][]int {
	if len(ranges) == 0 {
		return [][]int{{}}
	}

	var result [][]int
	var generate func(int, []int)
	generate = func(idx int, current []int) {
		if idx == len(ranges) {
			combo := make([]int, len(current))
			copy(combo, current)
			result = append(result, combo)
			return
		}

		for i := ranges[idx].min; i < ranges[idx].max; i++ {
			generate(idx+1, append(current, i))
		}
	}

	generate(0, []int{})
	return result
}
