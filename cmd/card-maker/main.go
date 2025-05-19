package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/wzyjerry/harmonies/data"
	"github.com/wzyjerry/harmonies/pkg/types"
	"golang.org/x/image/font"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var game types.GameData
var font18, font30 font.Face

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

	f, err := truetype.Parse(data.FontData)
	if err != nil {
		panic(err)
	}
	font18 = truetype.NewFace(f, &truetype.Options{
		Size: 18,
	})
	font30 = truetype.NewFace(f, &truetype.Options{
		Size: 30,
	})
}

const (
	w               = 300
	h               = 520
	size            = 40
	width           = 2 * size
	padding         = 2
	sizeWithPadding = size + padding
	horiz           = sizeWithPadding * 3 / 2
)

var (
	vert = math.Sqrt(3) * sizeWithPadding
)

func DrawTile() image.Image {
	dc := gg.NewContext(width, width)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.DrawRegularPolygon(6, width/2, width/2, size, 0)
	dc.SetRGBA255(241, 211, 161, 255)
	dc.Fill()
	return dc.Image()
}

func GetPatternSize(pattern []*types.Token) (minQ, maxQ, minR, maxR int) {
	minQ, maxQ = math.MaxInt, math.MinInt
	minR, maxR = math.MaxInt, math.MinInt
	for _, token := range pattern {
		if int(token.DeltaQ) < minQ {
			minQ = int(token.DeltaQ)
		}
		if int(token.DeltaQ) > maxQ {
			maxQ = int(token.DeltaQ)
		}
		if int(token.DeltaR) < minR {
			minR = int(token.DeltaR)
		}
		if int(token.DeltaR) > maxR {
			maxR = int(token.DeltaR)
		}
	}
	return
}

func DrawTileWithQR(dc *gg.Context, tile image.Image, q int, r int) {
	// dc.DrawImageAnchored(tile, w/2+q*horiz+30, h/2+int(float64(2*r+q)*vert/2)+180, 0.5, 0.5)
	dc.DrawImageAnchored(tile, w/2+q*horiz, h/2+int(float64(2*r+q)*vert/2), 0.5, 0.5)
}

func DrawToken(col color.Color, layer int) image.Image {
	const cx, cy = width / 2, width / 2
	dc := gg.NewContext(width, width)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	const r = size - 20 // 椭圆半径
	height := 15.0
	dc.Translate(0, 5-float64(layer)*height)

	// 底面边线
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	dc.DrawEllipse(cx, cy+height, r, r/2)
	dc.Stroke()

	// 侧面
	dc.SetColor(col)
	dc.DrawRectangle(cx-r, cy, r*2, height)
	dc.FillPreserve()
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	dc.Stroke()

	// 底面填充
	dc.SetColor(col)
	dc.DrawEllipse(cx, cy+height, r, r/2)
	dc.Fill()

	// 顶面
	dc.SetColor(col)
	dc.DrawEllipse(cx, cy, r, r/2)
	dc.FillPreserve()
	// 暗部
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 50})
	dc.FillPreserve()
	// 边线
	dc.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	dc.Stroke()

	return dc.Image()
}

func DrawScore(col color.Color, score int32) image.Image {
	dc := gg.NewContext(35, 70)
	animal := DrawAnimal(0)
	dc.DrawImageAnchored(animal, 35/2, 35/2, 0.5, 0.5)

	dc.SetColor(col)
	dc.DrawRectangle(5, 40, 25, 25)
	dc.FillPreserve()
	dc.SetRGBA255(255, 255, 255, 30)
	dc.Fill()

	dc.SetFontFace(font18)
	dc.SetRGBA255(0, 0, 0, 150)
	dc.DrawStringAnchored(fmt.Sprintf("%d", score), 35/2, 35+15, 0.5, 0.5)
	return dc.Image()
}

func DrawTitle(col color.Color, title string) image.Image {
	dc := gg.NewContext(w, 70)
	dc.SetColor(col)
	dc.SetFontFace(font30)
	dc.DrawStringAnchored(title, w/2, 35/2, 0.5, 0.5)
	return dc.Image()
}

func main() {
	colors := map[types.Color]color.RGBA{
		types.Color_ColorRed:    {R: 216, G: 40, B: 75, A: 255},
		types.Color_ColorGreen:  {R: 110, G: 182, B: 46, A: 255},
		types.Color_ColorBlue:   {R: 62, G: 126, B: 145, A: 255},
		types.Color_ColorYellow: {R: 251, G: 175, B: 37, A: 255},
		types.Color_ColorBrown:  {R: 111, G: 62, B: 59, A: 255},
		types.Color_ColorGray:   {R: 120, G: 103, B: 102, A: 255},
	}
	dc := gg.NewContext(w, h)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	animal := game.Animals[12]

	color := colors[types.Color_ColorBrown]
	switch animal.Kind {
	case types.POI_POIBuilding:
		color = colors[types.Color_ColorRed]
	case types.POI_POITree:
		color = colors[types.Color_ColorGreen]
	case types.POI_POIWater:
		color = colors[types.Color_ColorBlue]
	case types.POI_POIField:
		color = colors[types.Color_ColorYellow]
	case types.POI_POIMountain:
		color = colors[types.Color_ColorGray]
	}
	for i, score := range animal.Scores {
		dc.DrawImageAnchored(DrawScore(color, score), w-20, 40+70*(len(game.Animals[7].Scores)-i-1), 0.5, 0.5)
	}
	title := DrawTitle(color, animal.Name)
	dc.DrawImageAnchored(title, w/2, 45, 0.5, 0.5)

	dc.Translate(w/2, h/2)
	dc.Scale(1, 0.866)
	dc.Translate(-w/2, -h/2)

	tile := DrawTile()
	for _, t := range animal.Pattern {
		DrawTileWithQR(dc, tile, int(t.DeltaQ), int(t.DeltaR))
		switch t.Poi {
		case types.POI_POIWater:
			DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorBlue], 0), int(t.DeltaQ), int(t.DeltaR))
		case types.POI_POIField:
			DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorYellow], 0), int(t.DeltaQ), int(t.DeltaR))
		case types.POI_POIBuilding:
			for layer := range t.Height {
				DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorRed], int(layer)), int(t.DeltaQ), int(t.DeltaR))
			}
		case types.POI_POIMountain:
			for layer := range t.Height {
				DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorGray], int(layer)), int(t.DeltaQ), int(t.DeltaR))
			}
		case types.POI_POITree:
			for layer := range t.Height {
				if layer != t.Height-1 {
					DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorBrown], int(layer)), int(t.DeltaQ), int(t.DeltaR))
				} else {
					DrawTileWithQR(dc, DrawToken(colors[types.Color_ColorGreen], int(layer)), int(t.DeltaQ), int(t.DeltaR))
				}
			}
		}
		if t.Animal {
			DrawTileWithQR(dc, DrawAnimal(t.Height), int(t.DeltaQ), int(t.DeltaR))
		}
	}
	dc.SavePNG("../../output/card.png")
}

func DrawAnimal(layer int32) image.Image {
	dc := gg.NewContext(width, width)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.Translate(0, 10-float64(layer)*15)

	dc.DrawRegularPolygon(6, width/2, width/2, size/3, -math.Pi/2)
	dc.SetRGB255(252, 118, 52)
	dc.Fill()

	return dc.Image()
}

func PrintCard(card *types.Card) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(card.Name)
	buf.WriteString("]<")
	buf.WriteString(card.Kind.String())
	buf.WriteString(">\nscores: ")
	buf.WriteString(fmt.Sprintf("%v", card.Scores))
	buf.WriteString("\n")
	for _, pattern := range card.Pattern {
		buf.WriteString("  ")
		if pattern.Animal {
			buf.WriteString("Aanimal at ")
		}
		buf.WriteString("<")
		buf.WriteString(pattern.Poi.String())
		buf.WriteString(">")
		buf.WriteString(fmt.Sprintf(" [%d]", pattern.Height))
		if !pattern.Animal {
			buf.WriteString(fmt.Sprintf(" (%d, %d)\n", pattern.DeltaQ, pattern.DeltaR))
		} else {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
