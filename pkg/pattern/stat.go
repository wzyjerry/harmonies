package pattern

import (
	"bytes"
	"fmt"

	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Stat struct {
	Areas   []*Area
	Pattern *Pattern
}

type Area struct {
	POI   types.POI
	Hexes []HexInfo
}

type HexInfo struct {
	Hex               cube.Hex
	Distance          int
	NeighborTopColors int
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

func pop[T any](q *[]T) T {
	h := (*q)[0]
	*q = (*q)[1:]
	return h
}

func (p *Pattern) Stat() *Stat {
	start, ok := p.FindStarter()
	if !ok {
		return &Stat{}
	}
	p = p.Clone()
	stat := &Stat{
		Areas:   make([]*Area, 0),
		Pattern: p,
	}

	// 统计区块信息
	// 核心逻辑：选取一个网格，宽搜统计相邻的 POI 相同的网格形成区域 Area
	//         记录过程中见到的其余网格坐标作为其他区域候选起点
	visited := make(map[cube.Hex]Nil)

	seedQ := make([]cube.Hex, 0, 32)
	seedQ = append(seedQ, start)

	// 如果存在候选起点，选择一个作为当前区域种子
	for len(seedQ) > 0 {
		seed := pop(&seedQ)
		if _, ok := visited[seed]; ok {
			continue
		}
		poi := p.Get(seed).POI()

		// 当前区域队列
		areaQ := make([]HexInfo, 0, 32)
		hexes := make([]HexInfo, 0)
		areaQ = append(areaQ, HexInfo{
			Hex:      seed,
			Distance: 0,
		})
		for len(areaQ) > 0 {
			hexWithDistance := pop(&areaQ)
			hex := hexWithDistance.Hex
			visited[hex] = Nil{}
			topColor := make(map[types.Color]Nil)
			for dir := range cube.CubeDirectionCount {
				next := cube.CubeNeighbor(hex, dir)
				t := p.Get(next)
				if t == nil {
					continue
				}
				topColor[t.Top()] = Nil{}
				if _, ok := visited[next]; ok {
					continue
				}
				npoi := t.POI()
				if npoi == poi {
					areaQ = append(areaQ, HexInfo{
						Hex:      next,
						Distance: hexWithDistance.Distance + 1,
					})
				} else {
					seedQ = append(seedQ, next)
				}
			}
			hexes = append(hexes, HexInfo{
				Hex:               hex,
				Distance:          hexWithDistance.Distance,
				NeighborTopColors: len(topColor),
			})
		}
		stat.Areas = append(stat.Areas, &Area{
			POI:   poi,
			Hexes: hexes,
		})
	}
	return stat
}
