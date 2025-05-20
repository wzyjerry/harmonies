package finder

import (
	"container/heap"
	"fmt"
	"maps"

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
	perfabs := make(map[string]int)
	info := f.Cards[0]
	name := pattern.PoiHeight2Prefab(info.AnimalPOI, int32(info.AnimalHeight))
	perfab := pattern.Perfab2Func(name)
	perfabs[name]++
	animalHex := cube.NewHex(0, 0)
	perfab(seed, animalHex)
	usage := 1
	usageToken := pattern.Perfab2Tokens(name)
	distance := 0
	for _, tile := range info.PatternGroup[0] {
		name = pattern.PoiHeight2Prefab(tile.Poi, tile.Height)
		perfab = pattern.Perfab2Func(name)
		perfabs[name]++
		tileHex := cube.NewHex(int(tile.DeltaQ), int(tile.DeltaR))
		perfab(seed, tileHex)
		usage++
		usageToken += pattern.Perfab2Tokens(name)
		dis := tileHex.Distance()
		if dis > distance {
			distance = dis
		}
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	init := Scene{
		pattern:    seed,
		usage:      usage,
		usageToken: usageToken,
		perfabs:    perfabs,
		distance:   distance,
	}
	f.Solve(&init)
	heap.Push(pq, init)

	exists := make(map[string]Nil)
	exists[init.pattern.Hash()] = Nil{}

	found := false
	foundUsage := 0

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(Scene)
		pat := curr.pattern
		stat := pat.Stat()

		// 对于每个格子的空白邻居，尝试放置预制件
		// 如果放置后的情形未出现过，放置动物，入队
		for _, blocks := range stat.Areas {
			for _, hex := range blocks.Hexes {
				for dir := range cube.CubeDirectionCount {
					next := hex.Hex.Neighbor(dir)
					if next.Distance() > 4 {
						// 超出地图
						continue
					}
					if pat.Get(next) != nil {
						continue
					}

					// 格子是空的，尝试放置
					for name, perfab := range f.Perfabs {
						nextPattern := pat.Clone()
						perfab(nextPattern, next)
						hash := nextPattern.Hash()
						if _, ok := exists[hash]; ok {
							continue
						}
						exists[hash] = Nil{}

						// 预制件放置成功，尝试放置动物
						perfabs := maps.Clone(curr.perfabs)
						perfabs[name]++
						distance := next.Distance()
						if curr.distance > distance {
							distance = curr.distance
						}
						scene := Scene{
							pattern:    nextPattern,
							usage:      curr.usage + 1,
							usageToken: curr.usageToken + pattern.Perfab2Tokens(name),
							perfabs:    perfabs,
							distance:   distance,
						}
						f.Solve(&scene)

						if scene.animals == f.Total {
							// 找到一个解
							if !found {
								found = true
								foundUsage = scene.usageToken
							} else {
								if scene.usageToken < foundUsage {
									foundUsage = scene.usageToken
								}
							}
							fmt.Println(foundUsage)
							fmt.Println(scene.pattern.Hash())
						} else {
							if !found {
								heap.Push(pq, scene)
							} else {
								if scene.usageToken < foundUsage-1 {
									heap.Push(pq, scene)
								}
							}
						}
					}
				}
			}
		}
	}
}
