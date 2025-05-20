package pattern

import (
	"slices"

	"github.com/wzyjerry/harmonies/pkg/types"
)

type Tile struct {
	Height int
	Layers []types.Color
}

func NewTile() *Tile {
	return &Tile{
		Layers: make([]types.Color, 0, 3),
	}
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		Height: t.Height,
		Layers: slices.Clone(t.Layers),
	}
}

func (t *Tile) Top() types.Color {
	return t.Layers[t.Height]
}

func (t *Tile) POI() types.POI {
	switch t.Top() {
	case types.Color_ColorBlue:
		return types.POI_POIWater
	case types.Color_ColorGreen:
		return types.POI_POITree
	case types.Color_ColorGray:
		return types.POI_POIMountain
	case types.Color_ColorRed:
		return types.POI_POIBuilding
	case types.Color_ColorYellow:
		return types.POI_POIField
	}
	// 注意这里 types.Color_ColorBrown 是未定义 POI
	return types.POI_POIUnset
}
