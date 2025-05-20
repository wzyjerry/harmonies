// PlaceXX 对空白单元直接放置对应预制组件
package pattern

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

func PlaceWater(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBlue, 0)
}

func PlacePlain(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorYellow, 0)
}

func PlaceGress(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGreen, 0)
}

func PlaceTree(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorGreen, 1)
}

func PlaceForest(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorBrown, 1)
	p.Place(hex, types.Color_ColorGreen, 2)
}

func PlaceBuilding(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorRed, 0)
	p.Place(hex, types.Color_ColorRed, 1)
}

func PlaceRock(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorGray, 0)
}

func PlaceHill(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorBrown, 1)
}

func PlaceMountain(p *Pattern, hex cube.Hex) {
	p.Place(hex, types.Color_ColorBrown, 0)
	p.Place(hex, types.Color_ColorBrown, 1)
	p.Place(hex, types.Color_ColorBrown, 2)
}
