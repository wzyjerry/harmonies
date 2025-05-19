package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/wzyjerry/harmonies/data"
	"github.com/wzyjerry/harmonies/pkg/types"
	"golang.org/x/image/font"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var game types.GameData
var font14, font18, font30 font.Face

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
	font14 = truetype.NewFace(f, &truetype.Options{
		Size: 14,
	})
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
	vert   = math.Sqrt(3) * sizeWithPadding
	colors = map[types.Color]color.RGBA{
		types.Color_ColorRed:    {R: 216, G: 40, B: 75, A: 255},
		types.Color_ColorGreen:  {R: 110, G: 182, B: 46, A: 255},
		types.Color_ColorBlue:   {R: 62, G: 126, B: 145, A: 255},
		types.Color_ColorYellow: {R: 251, G: 175, B: 37, A: 255},
		types.Color_ColorBrown:  {R: 111, G: 62, B: 59, A: 255},
		types.Color_ColorGray:   {R: 120, G: 103, B: 102, A: 255},
	}
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

func DrawTileWithQR(dc *gg.Context, tile image.Image, q int, r int) {
	dc.DrawImageAnchored(tile, w/2+q*horiz, h/2+int(float64(2*r+q)*vert/2), 0.5, 0.5)
}

func CalcDrawRange(pattern []*types.Token) (width, height, left, top int32) {
	left, top = math.MaxInt32, math.MaxInt32
	var right, bottom int32 = math.MinInt32, math.MinInt32
	for _, t := range pattern {
		posX := t.DeltaQ * horiz
		posY := (2*t.DeltaR + t.DeltaQ) * int32(vert) / 2
		l := posX - size
		r := posX + size
		t := posY - int32(vert)/2
		b := posY + int32(vert)/2
		if l < left {
			left = l
		}
		if r > right {
			right = r
		}
		if t < top {
			top = t
		}
		if b > bottom {
			bottom = b
		}
	}
	return right - left, bottom - top, left, top
}

func DrawToken(col color.Color, layer int32) image.Image {
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

func DrawCardInfo(col color.Color, card *types.Card) image.Image {
	dc := gg.NewContext(w, 300)
	dc.SetColor(col)
	dc.SetFontFace(font14)
	info := strings.Split(PrintCard(card), "\n")
	for i, line := range info {
		dc.DrawString(line, 10, 30+20*(float64(i)+1))
	}
	return dc.Image()
}

func DrawTileWithQRAndPosition(dc *gg.Context, tile image.Image, q int32, r int32, left int32, top int32) {
	dc.DrawImageAnchored(tile, int(-left+q*horiz), int(-top)+int(float64(2*r+q)*vert/2), 0.5, 0.5)
}

func DrawPattern(animal *types.Card) image.Image {
	patternWidth, patternHeight, left, top := CalcDrawRange(animal.Pattern)

	dc := gg.NewContext(int(patternWidth), int(patternHeight))

	tile := DrawTile()
	for _, t := range animal.Pattern {
		DrawTileWithQRAndPosition(dc, tile, t.DeltaQ, t.DeltaR, left, top)
		switch t.Poi {
		case types.POI_POIWater:
			DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorBlue], 0), t.DeltaQ, t.DeltaR, left, top)
		case types.POI_POIField:
			DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorYellow], 0), t.DeltaQ, t.DeltaR, left, top)
		case types.POI_POIBuilding:
			for layer := range t.Height {
				DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorRed], layer), t.DeltaQ, t.DeltaR, left, top)
			}
		case types.POI_POIMountain:
			for layer := range t.Height {
				DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorGray], layer), t.DeltaQ, t.DeltaR, left, top)
			}
		case types.POI_POITree:
			for layer := range t.Height {
				if layer != t.Height-1 {
					DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorBrown], layer), t.DeltaQ, t.DeltaR, left, top)
				} else {
					DrawTileWithQRAndPosition(dc, DrawToken(colors[types.Color_ColorGreen], layer), t.DeltaQ, t.DeltaR, left, top)
				}
			}
		}
		if t.Animal {
			DrawTileWithQRAndPosition(dc, DrawAnimal(t.Height), t.DeltaQ, t.DeltaR, left, top)
		}
	}
	return dc.Image()
}

func main() {
	for i, animal := range game.Animals {
		dc := gg.NewContext(w, h)

		dc.SetRGB(1, 1, 1)
		dc.Clear()
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
			dc.DrawImageAnchored(DrawScore(color, score), w-20, 40+70*(len(animal.Scores)-i-1), 0.5, 0.5)
		}
		title := DrawTitle(color, animal.Name)
		dc.DrawImageAnchored(title, w/2, 45, 0.5, 0.5)
		cardInfo := DrawCardInfo(color, animal)
		dc.DrawImageAnchored(cardInfo, w/2, 200, 0.5, 0.5)

		dc.Translate(w/2, h/2)
		dc.Scale(1, 0.866)
		dc.Translate(-w/2, -h/2)
		pattern := DrawPattern(animal)
		dc.DrawImageAnchored(pattern, w/2, h-50, 0.5, 0.5)

		dc.SavePNG(fmt.Sprintf("../../output/card_%d.png", i))
	}
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
