package pattern

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/types"
)

const (
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

func (p *Pattern) calcDrawRange() (width, height, left, top int) {
	p.RLock()
	defer p.RUnlock()
	left, top = math.MaxInt, math.MaxInt
	var right, bottom int = math.MinInt, math.MinInt
	for t := range p.data {
		posX := t.Q * horiz
		posY := (2*t.R + t.Q) * int(vert) / 2
		l := posX - size
		r := posX + size
		t := posY - int(vert)/2
		b := posY + int(vert)/2
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

func DrawTile() image.Image {
	dc := gg.NewContext(width, width)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.DrawRegularPolygon(6, width/2, width/2, size, 0)
	dc.SetRGBA255(241, 211, 161, 255)
	dc.Fill()
	return dc.Image()
}

func DrawTileWithQRAndPosition(dc *gg.Context, tile image.Image, q int, r int, left int, top int) {
	dc.DrawImageAnchored(tile, int(-left+q*horiz), int(-top)+int(float64(2*r+q)*vert/2), 0.5, 0.5)
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

func (p *Pattern) Display() (*gg.Context, int, int) {
	patternWidth, patternHeight, left, top := p.calcDrawRange()
	dc := gg.NewContext(patternWidth, patternHeight)
	dc.Translate(float64(dc.Width())/2, float64(dc.Height())/2)
	dc.Scale(1, 0.866)
	dc.Translate(-float64(dc.Width())/2, -float64(dc.Height())/2)
	tile := DrawTile()
	p.RLock()
	defer p.RUnlock()
	for hex, info := range p.data {
		DrawTileWithQRAndPosition(dc, tile, hex.Q, hex.R, left, top)

		for i, col := range info.Layers {
			DrawTileWithQRAndPosition(dc, DrawToken(colors[col], i), hex.Q, hex.R, left, top)
		}
	}
	return dc, left, top
}

func DrawAnimal(layer int) image.Image {
	dc := gg.NewContext(width, width)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.Translate(0, 10-float64(layer)*15)

	dc.DrawRegularPolygon(6, width/2, width/2, size/3, -math.Pi/2)
	dc.SetRGB255(252, 118, 52)
	dc.Fill()

	return dc.Image()
}
func (p *Pattern) DisplayWithAnimals(animalsAt [][]cube.Hex) *gg.Context {
	dc, left, top := p.Display()
	for _, lst := range animalsAt {
		for _, hex := range lst {
			DrawTileWithQRAndPosition(dc, DrawAnimal(p.Get(hex).Height+1), hex.Q, hex.R, left, top)
		}
	}

	return dc
}
