package pattern

import (
	"cmp"
	"crypto/md5"
	"fmt"
	"slices"

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

type Cube struct {
	hex  cube.Hex
	tile *Tile
}

func (p *Pattern) Hash() string {
	// 计算棋盘 Hash，防止全局同型再现
	p.RLock()
	defer p.RUnlock()
	usage := make([]Cube, 0, len(p.data))
	for k, v := range p.data {
		usage = append(usage, Cube{k, v})
	}
	slices.SortStableFunc(usage, func(a, b Cube) int {
		if a.hex.Q == b.hex.Q {
			return cmp.Compare(a.hex.R, b.hex.R)
		}
		return cmp.Compare(a.hex.Q, b.hex.Q)
	})
	// 计算 Hash
	hash := md5.New()
	for _, v := range usage {
		hash.Write([]byte(fmt.Sprintf("%d_%d_%d", v.hex.Q, v.hex.R, v.tile.Height)))
		for _, col := range v.tile.Layers {
			hash.Write([]byte(fmt.Sprintf("_%d", col)))
		}
		hash.Write([]byte("-"))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
