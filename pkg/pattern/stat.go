package pattern

import (
	"bytes"
	"fmt"

	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Stat struct {
	Areas []*Area
}

type Area struct {
	POI   types.POI
	Hexes []cube.Hex
}

func (a *Area) String() string {
	var buf bytes.Buffer
	buf.WriteString("Area[")
	buf.WriteString(a.POI.String())
	buf.WriteString(fmt.Sprintf("] contains %d hex(es)", len(a.Hexes)))
	return buf.String()
}

func (s *Stat) String() string {
	if len(s.Areas) == 0 {
		return "Nil"
	}
	var buf bytes.Buffer
	for _, a := range s.Areas {
		buf.WriteString(a.String())
		buf.WriteString("\n")
	}
	return buf.String()
}

type Nil = struct{}

func (p *Pattern) Stat() *Stat {
	start, ok := p.FindStarter()
	if !ok {
		return &Stat{}
	}
	stat := &Stat{
		Areas: make([]*Area, 0),
	}

	// 统计区块信息
	// 核心逻辑：选取一个网格，宽搜统计相邻的 POI 相同的网格形成区域 Area
	//         记录过程中见到的其余网格坐标作为其他区域候选起点
	visited := make(map[cube.Hex]Nil)

	seedQ := make([]cube.Hex, 0, 128)
	seedQ = append(seedQ, start)

	pop := func(q *[]cube.Hex) cube.Hex {
		h := (*q)[0]
		*q = (*q)[1:]
		return h
	}

	// 如果存在候选起点，选择一个作为当前区域种子
	for len(seedQ) > 0 {
		seed := pop(&seedQ)
		if _, ok := visited[seed]; ok {
			continue
		}
		poi := p.Get(seed).POI()

		// 当前区域队列
		areaQ := make([]cube.Hex, 0, 128)
		hexes := make([]cube.Hex, 0)
		areaQ = append(areaQ, seed)
		for len(areaQ) > 0 {
			hex := pop(&areaQ)
			hexes = append(hexes, hex)
			visited[hex] = Nil{}
			for dir := range cube.CubeDirectionCount {
				next := cube.CubeNeighbor(hex, dir)
				t := p.Get(next)
				if t == nil {
					continue
				}
				if _, ok := visited[next]; ok {
					continue
				}
				npoi := t.POI()
				if npoi == poi {
					areaQ = append(areaQ, next)
				} else {
					seedQ = append(seedQ, next)
				}
			}
		}
		stat.Areas = append(stat.Areas, &Area{
			POI:   poi,
			Hexes: hexes,
		})
	}
	return stat
}
