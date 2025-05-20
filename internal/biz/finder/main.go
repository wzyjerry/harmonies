package finder

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/pattern"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Nil = struct{}

type Finder struct {
	Perfabs map[string]pattern.Perfab
	Cards   []*CardInfo
	Total   int
}

type CardInfo struct {
	// 除去 Animal 格后剩余方块
	// 全部 6 种旋转后的描述
	PatternGroup [][]*types.Token

	// 卡面完成所需放置的动物数量
	AnimalCount int

	// 动物所在的兴趣点
	AnimalPOI types.POI

	// 动物所在的层高
	AnimalHeight int

	Name string
}

func New(animals []*types.Card) *Finder {
	perfabs := make(map[string]pattern.Perfab)
	cards := make([]*CardInfo, 0, len(animals))
	total := 0
	for _, animal := range animals {
		info := &CardInfo{
			PatternGroup: make([][]*types.Token, cube.CubeDirectionCount),
			AnimalCount:  len(animal.Scores),
			Name:         animal.Name,
		}
		total += info.AnimalCount
		for _, tile := range animal.Pattern {
			name := pattern.PoiHeight2Prefab(tile.Poi, tile.Height)
			perfabs[name] = pattern.Perfab2Func(name)
			if tile.Animal {
				info.AnimalPOI = tile.Poi
				info.AnimalHeight = int(tile.Height)
				continue
			}
			hex := cube.NewHex(int(tile.DeltaQ), int(tile.DeltaR))
			info.PatternGroup[0] = append(info.PatternGroup[0], &types.Token{
				Poi:    tile.Poi,
				DeltaQ: int32(hex.Q),
				DeltaR: int32(hex.R),
				Height: tile.Height,
			})
			for i := 1; i < cube.CubeDirectionCount; i++ {
				hex = hex.RotateRight60()
				info.PatternGroup[i] = append(info.PatternGroup[i], &types.Token{
					Poi:    tile.Poi,
					DeltaQ: int32(hex.Q),
					DeltaR: int32(hex.R),
					Height: tile.Height,
				})
			}
		}
		cards = append(cards, info)
	}

	f := &Finder{
		Perfabs: perfabs,
		Cards:   cards,
		Total:   total,
	}
	return f
}
