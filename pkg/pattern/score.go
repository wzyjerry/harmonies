package pattern

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

// 计算不含河流的地形总分
// 计分规则如下：
//  1. 对于 POITree，如果 Top 是 Brown 记 0 分
//     否则根据高度记 1、3、7 分
//  2. 对于 POIMountain，如果 len(area) == 1 记 0 分
//     否则根据高度记 1、3、7 分
//  3. 对于 POIField，如果 len(area) == 1 记 0 分
//     否则记 5 分
//  4. 对于 POIBuilding，检查每个建筑的邻居的 6 个格子（如果存在）
//     如果顶部的颜色大于等于 3 种，记 5 分
//     否则记 0 分
func (s *Stat) TerrainVPWithoutWater() int {
	total := 0
	for _, area := range s.Areas {
		switch area.POI {
		case types.POI_POITree:
			for _, hex := range area.Hexes {
				t := s.Pattern.Get(hex.Hex)
				if t.Top() == types.Color_ColorBrown {
					continue
				}
				total += vpForHeight(t.Height)
			}
		case types.POI_POIMountain:
			if len(area.Hexes) == 1 {
				continue
			}
			for _, hex := range area.Hexes {
				total += vpForHeight(s.Pattern.Get(hex.Hex).Height)
			}
		case types.POI_POIField:
			if len(area.Hexes) == 1 {
				continue
			}
			total += 5
		case types.POI_POIBuilding:
			for _, hex := range area.Hexes {
				if hex.NeighborTopColors >= 3 {
					total += 5
				}
			}
		}
	}
	return total
}

func vpForHeight(height int) int {
	switch height {
	case 0:
		return 1
	case 1:
		return 3
	case 2:
		return 7
	}
	return 0
}

// 计算 A 面河流得分
// 查找最长河流，区域计算时已经记录了每个区域到种子的距离
// 这里只需要对于 POIWater 区域，找到区域内距离最大的点
// 作为起始点宽搜看所能到达的最远距离
// 记0、2、5、8、11、+4
func (s *Stat) WaterVPForA() int {
	var loooooooooooooooongest int
	for _, area := range s.Areas {
		if area.POI != types.POI_POIWater {
			continue
		}
		if len(area.Hexes) < 2 {
			continue
		}
		water := make(map[cube.Hex]Nil)
		max := area.Hexes[0].Distance
		maxHex := area.Hexes[0].Hex
		for _, hex := range area.Hexes {
			water[hex.Hex] = Nil{}
			if hex.Distance > max {
				max = hex.Distance
				maxHex = hex.Hex
			}
		}

		q := make([]HexInfo, 0, len(water))
		q = append(q, HexInfo{
			Hex:      maxHex,
			Distance: 0,
		})
		delete(water, maxHex)
		for len(q) > 0 {
			hex := pop(&q)
			for dir := range cube.CubeDirectionCount {
				next := hex.Hex.Neighbor(dir)
				if _, ok := water[next]; !ok {
					continue
				}
				delete(water, next)
				newDistance := hex.Distance + 1
				q = append(q, HexInfo{
					Hex:      next,
					Distance: newDistance,
				})
				if newDistance > loooooooooooooooongest {
					loooooooooooooooongest = newDistance
				}
			}
		}
	}
	if loooooooooooooooongest == 0 {
		return 0
	}
	total := 2
	if loooooooooooooooongest > 1 {
		total += 3 * (loooooooooooooooongest - 1)
	}
	if loooooooooooooooongest > 4 {
		total += loooooooooooooooongest - 4
	}
	return total
}
