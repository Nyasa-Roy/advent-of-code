package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const TWOPI = 2 * math.Pi

type Tile struct {
	X float64
	Y float64
}

func (t Tile) computeArea(other Tile) float64 {
	return (math.Abs(t.X-other.X) + 1) * (math.Abs(t.Y-other.Y) + 1)
}

func (t Tile) angle(curr, next *Tile) float64 {
	x1 := curr.X - t.X
	y1 := curr.Y - t.Y
	x2 := next.X - t.X
	y2 := next.Y - t.Y

	theta1 := math.Atan2(y1, x1)
	theta2 := math.Atan2(y2, x2)
	dtheta := theta2 - theta1

	for dtheta > math.Pi {
		dtheta -= TWOPI
	}
	for dtheta < -math.Pi {
		dtheta += TWOPI
	}
	return dtheta
}

func (t Tile) pointOnAxisAlignedSegment(a, b *Tile) bool {
	if a.X == b.X && t.X == a.X {
		if a.Y < b.Y {
			return t.Y >= a.Y && t.Y <= b.Y
		}
		return t.Y >= b.Y && t.Y <= a.Y
	}
	if a.Y == b.Y && t.Y == a.Y {
		if a.X < b.X {
			return t.X >= a.X && t.X <= b.X
		}
		return t.X >= b.X && t.X <= a.X
	}
	return false
}

type Node struct {
	Tile *Tile
	Next *Node
}

func NewNode(tile *Tile) *Node {
	return &Node{Tile: tile}
}

type rectBounds struct {
	xLo float64
	xHi float64
	yLo float64
	yHi float64
}

func newRectBounds(a, b Tile) rectBounds {
	xLo, xHi := a.X, b.X
	if xLo > xHi {
		xLo, xHi = xHi, xLo
	}
	yLo, yHi := a.Y, b.Y
	if yLo > yHi {
		yLo, yHi = yHi, yLo
	}
	return rectBounds{xLo, xHi, yLo, yHi}
}

func (r rectBounds) contains(t *Tile) bool {
	return t.X > r.xLo && t.X < r.xHi && t.Y > r.yLo && t.Y < r.yHi
}

func (r rectBounds) isCrossedHorizontally(curr, next *Tile) bool {
	if curr.Y <= r.yLo || curr.Y >= r.yHi {
		return false
	}
	return (curr.X >= r.xHi && next.X <= r.xLo) || (curr.X <= r.xLo && next.X >= r.xHi)
}

func (r rectBounds) isCrossedVertically(curr, next *Tile) bool {
	if curr.X <= r.xLo || curr.X >= r.xHi {
		return false
	}
	return (curr.Y >= r.yHi && next.Y <= r.yLo) || (curr.Y <= r.yLo && next.Y >= r.yHi)
}

func (r rectBounds) edgeCrossesOrEnters(curr, next *Tile) bool {
	if r.contains(curr) {
		return true
	}
	if r.isCrossedHorizontally(curr, next) {
		return true
	}
	return r.isCrossedVertically(curr, next)
}

func cornerInside(corner Tile, start *Node) bool {
	const angleEps = 1e-7
	angle := 0.0
	onEdge := false

	update := func(curr, next *Tile) {
		if onEdge {
			return
		}
		if corner.pointOnAxisAlignedSegment(curr, next) {
			onEdge = true
			return
		}
		angle += corner.angle(curr, next)
	}

	for curr := start; ; curr = curr.Next {
		update(curr.Tile, curr.Next.Tile)
		if curr.Next == start {
			break
		}
	}
	return onEdge || math.Abs(angle) > (TWOPI-angleEps)
}

func (n *Node) formsValidRectangleWith(other *Node) bool {
	if n.Next == other {
		return true
	}

	bounds := newRectBounds(*n.Tile, *other.Tile)
	a := Tile{X: n.Tile.X, Y: other.Tile.Y}
	b := Tile{X: other.Tile.X, Y: n.Tile.Y}

	for curr := n; ; curr = curr.Next {
		if curr != n && curr != other && curr.Next != other {
			if bounds.edgeCrossesOrEnters(curr.Tile, curr.Next.Tile) {
				return false
			}
		}
		if curr.Next == n {
			break
		}
	}
	return cornerInside(a, n) && cornerInside(b, n)
}

func toTilePointers(points []Tile) []*Tile {
	ptrs := make([]*Tile, len(points))
	for i := range points {
		ptrs[i] = &points[i]
	}
	return ptrs
}

func createNodeRing(pointPtrs []*Tile) []*Node {
	nodes := make([]*Node, len(pointPtrs))
	for i, curr := range pointPtrs {
		nodes[i] = NewNode(curr)
		if i > 0 {
			nodes[i-1].Next = nodes[i]
		}
	}
	nodes[len(nodes)-1].Next = nodes[0]
	return nodes
}

func findLargestRectangle(nodes []*Node, constrain bool) float64 {
	best := 0.0
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			area := nodes[i].Tile.computeArea(*nodes[j].Tile)
			if area > best {
				if !constrain || nodes[i].formsValidRectangleWith(nodes[j]) {
					best = area
				}
			}
		}
	}
	return best
}

func readPointsFromFile(path string) ([]Tile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	coords := make([]Tile, 0)
	for scanner.Scan() {
		var c Tile
		_, err := fmt.Sscanf(scanner.Text(), "%f,%f", &c.X, &c.Y)
		if err != nil {
			return nil, fmt.Errorf("error parsing line %q: %w", scanner.Text(), err)
		}
		coords = append(coords, c)
	}
	return coords, scanner.Err()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <inputfile> <true|false>")
		return
	}

	points, err := readPointsFromFile(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	constrain := os.Args[2] == "true"
	nodes := createNodeRing(toTilePointers(points))
	area := findLargestRectangle(nodes, constrain)

	fmt.Printf("Largest rectangle area: %d\n", int(area))
}
