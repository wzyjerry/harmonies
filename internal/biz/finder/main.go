package finder

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/pattern"
	"github.com/wzyjerry/harmonies/pkg/types"
)

type Nil = struct{}

type Finder struct {
	Perfabs []pattern.Perfab
	Cards   []*CardInfo
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
	AnimalHeight int32
}

func New(animals []*types.Card) *Finder {
	perfab := make(map[string]Nil)
	cards := make([]*CardInfo, 0, len(animals))
	for _, animal := range animals {
		info := &CardInfo{
			PatternGroup: make([][]*types.Token, cube.CubeDirectionCount),
			AnimalCount:  len(animal.Scores),
		}
		for _, tile := range animal.Pattern {
			perfab[pattern.PoiHeight2Prefab(tile.Poi, tile.Height)] = Nil{}
			if tile.Animal {
				info.AnimalPOI = tile.Poi
				info.AnimalHeight = tile.Height
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
	perfabs := make([]pattern.Perfab, 0, len(perfab))
	for k := range perfab {
		perfabs = append(perfabs, pattern.Perfab2Func(k))
	}

	f := &Finder{
		Perfabs: perfabs,
		Cards:   cards,
	}
	return f
}
