package pattern

import (
	"slices"

	"github.com/wzyjerry/harmonies/pkg/types"
)

type Tile struct {
	height int
	layers []types.Color
}

func NewTile() *Tile {
	return &Tile{
		layers: make([]types.Color, 0, 3),
	}
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		height: t.height,
		layers: slices.Clone(t.layers),
	}
}

func (t *Tile) Top() types.Color {
	return t.layers[t.height]
}

func (t *Tile) POI() types.POI {
	switch t.Top() {
	case types.Color_ColorBlue:
		return types.POI_POIWater
	case types.Color_ColorGreen,
		types.Color_ColorBrown:
		return types.POI_POITree
	case types.Color_ColorGray:
		return types.POI_POIMountain
	case types.Color_ColorRed:
		return types.POI_POIBuilding
	case types.Color_ColorYellow:
		return types.POI_POIField
	}
	return types.POI_POIUnset
}
