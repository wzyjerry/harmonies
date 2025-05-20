package cube

const (
	CubeDirectionRightDown = iota
	CubeDirectionRightUp
	CubeDirectionUp
	CubeDirectionLeftUp
	CubeDirectionLeftDown
	CubeDirectionDown
	CubeDirectionCount
)

type Hex struct {
	Q int
	R int
}

func NewHex(q int, r int) Hex {
	return Hex{Q: q, R: r}
}

var axialDirectionVectors = [][2]int{
	{+1, 0}, {+1, -1}, {0, -1},
	{-1, 0}, {-1, +1}, {0, +1},
}

func CubeDirection(direction int) [2]int {
	return axialDirectionVectors[direction]
}

func CubeAdd(hex Hex, vec [2]int) Hex {
	return NewHex(hex.Q+vec[0], hex.R+vec[1])
}

func CubeNeighbor(hex Hex, direction int) Hex {
	return CubeAdd(hex, CubeDirection(direction))
}
