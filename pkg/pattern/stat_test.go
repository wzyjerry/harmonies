package pattern_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wzyjerry/harmonies/data"
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/pattern"
	"github.com/wzyjerry/harmonies/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var game types.GameData

func yaml2json(b []byte) ([]byte, error) {
	var data map[string]any
	err := yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func init() {
	raw, err := yaml2json(data.Data)
	if err != nil {
		panic(err)
	}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(raw, &game)
	if err != nil {
		panic(err)
	}
}

func TestStat(t *testing.T) {
	for i, animal := range game.Animals {
		p := pattern.NewPattern()
		for _, token := range animal.Pattern {
			hex := cube.NewHex(int(token.DeltaQ), int(token.DeltaR))
			switch token.Poi {
			case types.POI_POIField:
				assert.True(t, p.CanPlace(hex, types.Color_ColorYellow, 0))
				p.Place(hex, types.Color_ColorYellow, 0)
			case types.POI_POIWater:
				assert.True(t, p.CanPlace(hex, types.Color_ColorBlue, 0))
				p.Place(hex, types.Color_ColorBlue, 0)
			case types.POI_POIBuilding:
				for layer := range token.Height {
					assert.True(t, p.CanPlace(hex, types.Color_ColorRed, int(layer)))
					p.Place(hex, types.Color_ColorRed, int(layer))
				}
			case types.POI_POIMountain:
				for layer := range token.Height {
					assert.True(t, p.CanPlace(hex, types.Color_ColorGray, int(layer)))
					p.Place(hex, types.Color_ColorGray, int(layer))
				}
			case types.POI_POITree:
				for layer := range token.Height - 1 {
					assert.True(t, p.CanPlace(hex, types.Color_ColorBrown, int(layer)))
					p.Place(hex, types.Color_ColorBrown, int(layer))
				}
				assert.True(t, p.CanPlace(hex, types.Color_ColorGreen, int(token.Height)-1))
				p.Place(hex, types.Color_ColorGreen, int(token.Height)-1)
			}
		}
		fmt.Println("animal_", i, animal.Name)
		stat := p.Stat()
		fmt.Println(stat)
		fmt.Println("VPWithoutWater", stat.TerrainVPWithoutWater())
		fmt.Println("VPAWater", stat.WaterVPForA())
		fmt.Println()
	}
}
