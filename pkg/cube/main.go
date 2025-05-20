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

func (h Hex) Add(vec [2]int) Hex {
	return NewHex(h.Q+vec[0], h.R+vec[1])
}
