package cube

func (h Hex) Subtract(b Hex) Hex {
	return NewHex(h.Q-b.Q, h.R-b.R)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (h Hex) DistanceBy(b Hex) int {
	vec := h.Subtract(b)
	return vec.Distance()
}

func (h Hex) Distance() int {
	s := -h.Q - h.R
	return (abs(h.Q) + abs(h.R) + abs(s)) / 2
}
