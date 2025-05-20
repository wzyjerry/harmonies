package finder

import (
	"container/heap"
	"fmt"

	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/pattern"
)

// 填充 Perfab，Solve，最终目标使 scene 放置的动物数量与 Total 相同
func (f *Finder) Search() {
	// 没有动物，直接返回
	if len(f.Cards) == 0 {
		return
	}

	// 构建起始种子：第一张动物卡第一种形态
	seed := pattern.NewPattern()
	visited := make(map[cube.Hex]Nil)
	perfabs := make(map[string]int)
	info := f.Cards[0]
	name := pattern.PoiHeight2Prefab(info.AnimalPOI, int32(info.AnimalHeight))
	perfab := pattern.Perfab2Func(name)
	perfabs[name]++
	animalHex := cube.NewHex(0, 0)
	perfab(seed, animalHex)
	visited[animalHex] = Nil{}
	usage := 1
	distance := 0
	for _, tile := range info.PatternGroup[0] {
		name = pattern.PoiHeight2Prefab(tile.Poi, tile.Height)
		perfab = pattern.Perfab2Func(name)
		perfabs[name]++
		tileHex := cube.NewHex(int(tile.DeltaQ), int(tile.DeltaR))
		perfab(seed, tileHex)
		visited[tileHex] = Nil{}
		usage++
		dis := tileHex.Distance()
		if dis > distance {
			distance = dis
		}
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	init := Scene{
		visited:  visited,
		pattern:  seed,
		usage:    usage,
		perfabs:  perfabs,
		distance: distance,
	}
	f.Solve(&init)
	heap.Push(pq, init)
	fmt.Println(init.pattern.Hash())
}
