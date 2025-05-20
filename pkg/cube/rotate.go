package cube

func (h Hex) RotateRight60() Hex {
	r, s := h.R, -h.Q-h.R
	return NewHex(-r, -s)
}

func (h Hex) RotateLeft60() Hex {
	q, s := h.Q, -h.Q-h.R
	return NewHex(-s, -q)
}
