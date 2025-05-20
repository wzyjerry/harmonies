package pattern

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
)

func (p *Pattern) Clone() *Pattern {
	p.RLock()
	defer p.RUnlock()
	data := make(map[cube.Hex]*Tile)
	for k, v := range p.data {
		data[k] = v.Clone()
	}
	return &Pattern{
		data: data,
	}
}

func (p *Pattern) Get(hex cube.Hex) *Tile {
	p.RLock()
	defer p.RUnlock()
	return p.data[hex]
}

func (p *Pattern) FindStarter() (cube.Hex, bool) {
	p.RLock()
	defer p.RUnlock()
	for k := range p.data {
		return k, true
	}
	return cube.NewHex(0, 0), false
}
