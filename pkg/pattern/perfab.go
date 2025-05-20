// PerfabXX 对空白单元直接放置对应预制组件
package pattern

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Perfab func(p *Pattern, hex cube.Hex)

func PerfabNothing(p *Pattern, hex cube.Hex) {}

func PerfabWater(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBlue, 0)
}

func PerfabPlain(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorYellow, 0)
}

func PerfabGrass(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGreen, 0)
}

func PerfabTree(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorGreen, 1)
}

func PerfabForest(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorBrown, 1)
	p.Place(hex, types.Color_ColorGreen, 2)
}

func PerfabBuilding(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorRed, 0)
	p.Place(hex, types.Color_ColorRed, 1)
}

func PerfabRock(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGray, 0)
}

func PerfabHill(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGray, 0)
	p.Place(hex, types.Color_ColorGray, 1)
}

func PerfabMountain(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGray, 0)
	p.Place(hex, types.Color_ColorGray, 1)
	p.Place(hex, types.Color_ColorGray, 2)
}

func PoiHeight2Prefab(poi types.POI, height int32) string {
	switch poi {
	case types.POI_POIWater:
		if height == 1 {
			return "water"
		}
		return ""
	case types.POI_POIField:
		if height == 1 {
			return "plain"
		}
		return ""
	case types.POI_POIBuilding:
		{
			if height == 2 {
				return "building"
			}
			return ""
		}
	case types.POI_POITree:
		if height == 1 {
			return "grass"
		}
		if height == 3 {
			return "forest"
		}
		if height == 2 {
			return "tree"
		}
		return ""
	case types.POI_POIMountain:
		if height == 1 {
			return "rock"
		}
		if height == 2 {
			return "hill"
		}
		if height == 3 {
			return "mountain"
		}
		return ""
	}
	return ""
}

func Perfab2Func(prefab string) Perfab {
	switch prefab {
	case "water":
		return PerfabWater
	case "plain":
		return PerfabPlain
	case "grass":
		return PerfabGrass
	case "tree":
		return PerfabTree
	case "forest":
		return PerfabForest
	case "building":
		return PerfabBuilding
	case "rock":
		return PerfabRock
	case "hill":
		return PerfabHill
	case "mountain":
		return PerfabMountain
	}
	return PerfabNothing
}

func Perfab2Tokens(prefab string) int {
	switch prefab {
	case "water":
		return 1
	case "plain":
		return 1
	case "grass":
		return 1
	case "tree":
		return 2
	case "forest":
		return 3
	case "building":
		return 2
	case "rock":
		return 1
	case "hill":
		return 2
	case "mountain":
		return 3
	}
	return 0
}
