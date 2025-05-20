package cube

var axialDirectionVectors = [][2]int{
	{+1, 0}, {+1, -1}, {0, -1},
	{-1, 0}, {-1, +1}, {0, +1},
}

func CubeDirection(direction int) [2]int {
	return axialDirectionVectors[direction]
}

func (h Hex) Neighbor(direction int) Hex {
	return h.Add(CubeDirection(direction))
}
