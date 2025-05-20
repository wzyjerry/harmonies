package pattern

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

// 放置规则
// 1. 0 层允许放 Gray、Red、Brown、Green、Blue、Yellow
// 2. 如果是 Gray，允许放在 1、2 层，要求最上层为 Gray
// 3. 如果是 Red，允许放在 1 层，要求最上层为 Gray、Red、Brown
// 4. 如果是 Brown，允许放在 1 层，要求最上层为 Brown
// 5. 如果是 Green，允许放在 1、2 层，要求最上层为 Brown
// 6. 拒绝其他情况
// CanPlace 不检查相邻性，只检查当前颜色是否允许放置在指定位置
func (p *Pattern) CanPlace(hex cube.Hex, color types.Color, layer int) bool {
	if layer < 0 || layer > 2 {
		return false
	}
	tile := p.Get(hex)
	if tile == nil {
		// 检查放置层
		if layer != 0 {
			return false
		}
		switch color {
		case types.Color_ColorGray,
			types.Color_ColorRed,
			types.Color_ColorBrown,
			types.Color_ColorGreen,
			types.Color_ColorBlue,
			types.Color_ColorYellow:
			return true
		default:
			return false
		}
	}
	// 检查放置层
	if layer != tile.height+1 {
		return false
	}
	top := tile.Top()
	switch color {
	case types.Color_ColorGray:
		if top == types.Color_ColorGray {
			return true
		}
		return false
	case types.Color_ColorRed:
		if layer != 1 {
			return false
		}
		if top == types.Color_ColorGray ||
			top == types.Color_ColorRed ||
			top == types.Color_ColorBrown {
			return true
		}
		return false
	case types.Color_ColorBrown:
		if layer != 1 {
			return false
		}
		if top == types.Color_ColorBrown {
			return true
		}
		return false
	case types.Color_ColorGreen:
		if top == types.Color_ColorBrown {
			return true
		}
		return false
	}
	return false
}
