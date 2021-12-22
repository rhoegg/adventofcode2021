package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Cuboid struct {
	X, Y, Z [2]int
}

func (c Cuboid) String() string {
return fmt.Sprintf("(%d..%d),(%d..%d)(%d..%d)", c.X[0], c.X[1], c.Y[0], c.Y[1], c.Z[0], c.Z[1])
}

func (c Cuboid) Touches(other Cuboid) bool {
	return c.TouchesX(other) && c.TouchesY(other) && c.TouchesZ(other)
}

func (c Cuboid) TouchesX(other Cuboid) bool {
	return ! (other.X[1] < c.X[0] || other.X[0] > c.X[1])
}

func (c Cuboid) TouchesY(other Cuboid) bool {
	return ! (other.Y[1] < c.Y[0] || other.Y[0] > c.Y[1])
}

func (c Cuboid) TouchesZ(other Cuboid) bool {
	return ! (other.Z[1] < c.Z[0] || other.Z[0] > c.Z[1])
}

func (c Cuboid) Inside(other Cuboid) bool {
	return c.InsideX(other) && c.InsideY(other) && c.InsideZ(other)
}

func (c Cuboid) InsideX(other Cuboid) bool {
	return c.X[0] >= other.X[0] && c.X[1] <= other.X[1]
}

func (c Cuboid) InsideY(other Cuboid) bool {
	return c.Y[0] >= other.Y[0] && c.Y[1] <= other.Y[1]
}

func (c Cuboid) InsideZ(other Cuboid) bool {
	return c.Z[0] >= other.Z[0] && c.Z[1] <= other.Z[1]
}

func (c Cuboid) SplitWithin(domain [2]int, other [2]int) [][2]int {
	if other[0] <= domain[0] && other[1] >= domain[1] {
		return [][2]int{domain}
	}
	if other[0] > domain[0] && other[1] < domain[1] {
		// three segments
		return [][2]int{{domain[0], other[0] - 1}, other, {other[1] + 1, domain[1]}}
	}

	if other[0] <= domain[0] {
		var segments [][2]int
		segments = append(segments, [2]int{domain[0], other[1]})
		segments = append(segments, [2]int{other[1] + 1, domain[1]})
		return segments
	} else {
		var segments [][2]int
		segments = append(segments, [2]int{domain[0], other[0] - 1})
		segments = append(segments, [2]int{other[0], domain[1]})
		return segments
	}
}

func (c Cuboid) CubesLeftOn(steps []InitStep) int64 {
	volume := int64(c.X[1] - c.X[0] + 1) *
		int64(c.Y[1] - c.Y[0] + 1) *
		int64(c.Z[1] - c.Z[0] + 1)
	if volume < 262144 {
		count := 0
		for x := c.X[0]; x <= c.X[1]; x++ {
			for y := c.Y[0]; y <= c.Y[1]; y++ {
				for z := c.Z[0]; z <= c.Z[1]; z++ {
					on := false
					for _, step := range steps {
						on = step.Apply(x, y, z, on)
					}
					if on {
						count++
					}
				}
			}
		}
		fmt.Printf("Small block done %v\n", c)
		return int64(count)
	}
	on := false
	for  _, step := range steps {
		if step.Cuboid.Touches(c) {
			if c.Inside(step.Cuboid) {
				on = step.Action
			} else {
				// cuboid is not uniform, split  and sum
				var fragments []Cuboid
				for _, splitX := range c.SplitWithin(c.X, step.X) {
					for _, splitY := range c.SplitWithin(c.Y, step.Y) {
						for _, splitZ := range c.SplitWithin(c.Z, step.Z) {
							fragments = append(fragments, Cuboid{X: splitX, Y: splitY, Z: splitZ})
						}
					}
				}

				sum := int64(0)
				//fmt.Printf("Down %v %v %v\n", fragments[0], c, step.Cuboid)
				for _, fragment := range fragments {
					sum += fragment.CubesLeftOn(steps)
				}
				//fmt.Println("Up")
				return sum
			}
		}
	}
	fmt.Printf("Solid block %v", c)

	if on {
		fmt.Println(" on")
		return volume
	}
	fmt.Println(" off")
	return 0
}

type InitStep struct {
	Cuboid
	Action bool
}


func (s *InitStep) Apply(x, y, z int, prev bool) bool {
	if x < s.X[0] || x > s.X[1] {
		return prev
	}
	if y < s.Y[0] || y > s.Y[1] {
		return prev
	}
	if z < s.Z[0] || z > s.Z[1] {
		return prev
	}
	return s.Action
}

func ParseInitStep(procedure string) InitStep {
	tokens := strings.Split(procedure, " ")
	action := tokens[0] == "on"
	dimensionsInput := strings.Split(tokens[1], ",")

	var dims [][2]int
	for _, d := range dimensionsInput {
		bounds := strings.Split(d[2:], "..")
		low, _ := strconv.Atoi(bounds[0])
		high, _ := strconv.Atoi(bounds[1])
		dims = append(dims, [2]int{low, high})
	}
	return InitStep{
		Cuboid: Cuboid{
			X: dims[0],
			Y: dims[1],
			Z: dims[2],
		},
		Action: action,
	}
}

func Bounds(steps []InitStep) Cuboid {
	bounds := Cuboid{}
	for _, step := range steps {
		if step.X[0] < bounds.X[0] {
			bounds.X[0] = step.X[0]
		}
		if step.X[1] > bounds.X[1] {
			bounds.X[1] = step.X[1]
		}
		if step.Y[0] < bounds.Y[0] {
			bounds.Y[0] = step.Y[0]
		}
		if step.Y[1] > bounds.Y[1] {
			bounds.Y[1] = step.Y[1]
		}
		if step.Z[0] < bounds.Z[0] {
			bounds.Z[0] = step.Z[0]
		}
		if step.Z[1] > bounds.Z[1] {
			bounds.Z[1] = step.Z[1]
		}
	}
	return bounds
}


func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	var initSteps []InitStep

	for _, step := range strings.Split(string(data), "\n") {
		initSteps = append(initSteps, ParseInitStep(step))
	}
	bounds := Bounds(initSteps)
	fmt.Println(bounds.CubesLeftOn(initSteps))
}
