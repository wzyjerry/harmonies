package finder

import (
	"cmp"
	"slices"

	"github.com/wzyjerry/harmonies/pkg/cube"
)

// 放置动物
// 由于动物所在地形存在重复，因此放置需要考虑所有选点
// 此处贪心直觉上成立：
// 优先采用无征用的选点，有共用选点时优先让满足总数少的动物使用
// 因为最终状态所有动物都有选点，因此分配不会影响正解
//
// 查找方式：对于每张动物优先查找潜在的动物选点
// 如果周围地形与某种旋转规则匹配，记录选点
// 使用上述贪心确定最终选点
//
// Solve 方法最终会调用一次 scene.CalcScore 更新得分信息
func (f *Finder) Solve(scene *Scene) {
	// 计算场景区块信息
	pattern := scene.pattern
	stat := pattern.Stat()

	points := make([]*AnimalPoint, len(f.Cards))
	// 对于每张动物卡
	for i, info := range f.Cards {
		potential := make(map[cube.Hex]Nil)
		// 选取匹配的区块
		for _, block := range stat.Areas {
			if block.POI != info.AnimalPOI {
				continue
			}

			// 查找潜在动物选点
			for _, hex := range block.Hexes {
				if hex.Height != info.AnimalHeight-1 {
					continue
				}
				potential[hex.Hex] = Nil{}
			}
		}

		// 地形匹配确认潜在选点
		points[i] = &AnimalPoint{
			animalCount: info.AnimalCount,
			animalID:    i,
			points:      make(map[cube.Hex]Nil),
		}
		for hex := range potential {
			apply := false
		outer:
			for _, group := range info.PatternGroup {
				for _, token := range group {
					aim := hex.Add([2]int{int(token.DeltaQ), int(token.DeltaR)})
					tile := pattern.Get(aim)
					if tile == nil {
						continue outer
					}
					if tile.POI() != token.Poi || tile.Height != int(token.Height)-1 {
						continue outer
					}
				}
				apply = true
				break
			}
			if apply {
				points[i].points[hex] = Nil{}
			}
		}
	}

	// 按选点数量排序，最终分配
	slices.SortStableFunc(points, func(a, b *AnimalPoint) int {
		return cmp.Compare(len(a.points), len(b.points))
	})

	pointerMap := make(map[cube.Hex]int)
	for _, point := range points {
		for hex := range point.points {
			pointerMap[hex]++
		}
	}

	total := 0
	animalsAt := make([][]cube.Hex, len(f.Cards))
	for _, point := range points {
		lst := make([]PointCount, 0, len(point.points))
		for hex := range point.points {
			count, ok := pointerMap[hex]
			if !ok {
				continue
			}
			lst = append(lst, PointCount{
				point: hex,
				count: count,
			})
		}
		slices.SortStableFunc(lst, func(a, b PointCount) int {
			return cmp.Compare(a.count, b.count)
		})
		place := min(len(lst), point.animalCount)
		animalsAt[point.animalID] = make([]cube.Hex, 0, place)
		for i := range place {
			animalsAt[point.animalID] = append(animalsAt[point.animalID], lst[i].point)
			delete(pointerMap, lst[i].point)
		}
		total += place
	}
	scene.animals = total
	scene.animalsAt = animalsAt
	scene.CalcScore()
}

type AnimalPoint struct {
	animalCount int
	animalID    int
	points      map[cube.Hex]Nil
}

type PointCount struct {
	point cube.Hex
	count int
}
