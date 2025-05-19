package data

import (
	_ "embed"
)

//go:embed JetBrainsMono-Regular-2.ttf
var FontData []byte

//go:embed animals.yaml
var Data []byte
