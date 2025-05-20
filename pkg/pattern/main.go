// pattern 表示一个连续的六边型区域，对长宽并无限制
package pattern

import (
	"sync"

	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Pattern struct {
	sync.RWMutex
	data map[cube.Hex]*Tile
}

func NewPattern() *Pattern {
	return &Pattern{
		data: make(map[cube.Hex]*Tile),
	}
}

// Place 放置颜色，要求调用前通过 CanPlace 检查
func (p *Pattern) Place(hex cube.Hex, color types.Color, layer int) {
	tile := p.Get(hex)
	if tile == nil {
		tile = NewTile()
		p.Lock()
		defer p.Unlock()
		p.data[hex] = tile
	}
	tile.Layers = append(tile.Layers, color)
	tile.Height = layer
}
