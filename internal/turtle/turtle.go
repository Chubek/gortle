package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Wrapping int
type PenMode int

const (
	WrappingWrap   Wrapping = iota
	WrappingFence  Wrapping
	WrappingWindow Wrapping

	PenPaint   PenMode = iota
	PenErase   PenMode
	PenReverse PenMode
)

var (
	WindowWidth  = 800
	WindowHeight = 600
)

type Point struct {
	X, Y float64
}

type Color struct {
	R, G, B, A uint8
}

type Turtle struct {
	x, y       float64
	angle      float64
	penDown    bool
	showTurtle bool
	recordPath bool
	wrapMode   Wrapping
	penMode    PenMode
	bgColor    Color
	fgColor    Color
	scale      float64
	minX, minY int32
	maxX, maxY int32
	spriteW    int32
	spriteH    int32
	fontSize   uint
	fontPath   string
	path       []Point
	renderer   *sdl.Renderer
	sprite     *sdl.Texture
	font       *ttf.Font
}

func (c Color) toSDLColor() {
	return sdl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

func (c Color) getFields() (uint8, uint8, uint8, uint8) {
	return c.r, c.g, c.b, c.a
}

func (c Color) getInverseFields() (uint8, uint8, uint8) {
	return 255 - c.r, 255 - c.g, 255 - c.b, c.a
}

func (p Point) toSDLPoint() {
	return sdl.Point{X: p.X, Y: p.Y}
}

func NewTurtle(r *sdl.Renderer, s *sdl.Texture) *Turtle {
	t := &Turtle{
		x:          0,
		y:          0,
		angle:      0,
		penDown:    false,
		showTurtle: false,
		recordPath: false,
		penMode:    PenPaint,
		wrapMode:   WrappingWrap,
		bgColor:    Color{255, 255, 255, 255},
		fgColor:    Color{0, 0, 0, 0},
		scale:      1.0,
		minX:       0,
		minY:       0,
		maxX:       WindowWidth - 1,
		maxY:       WindowHeight - 1,
		spriteW:    -1,
		spriteH:    -1,
		fontSize:   12,
		fontPath:   os.GetEnv("GORTLE_DEFAULT_FONTPATH"),
		path:       make([]Point, 0, 1024),
		renderer:   r,
		sprite:     s,
		font:       nil,
	}
	return t
}

func (t *Turtle) LoadTurtleImage(path string) error {
	if tex, err := img.LoadTexture(t.renderer, path); err != nil {
		return fmt.Errorf("turtleimage: img.LoadTexture failed: %v", err)
	}

	var query sdl.TextureInfo
	if err := tex.Query(&query); err != nil {
		tex.Destroy()
		return fmt.Errorf("turtleimage: tex.Query failed: %v", err)
	}

	t.sprite = tex
	t.spriteW = query.Width
	t.spriteH = query.Height

	return nil
}

func (t *Turtle) LoadFont() error {
	t.CloseFont()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("setlabelfont: ttf.Init failed: %v", err)
	}

	if f, err := ttf.OpenFont(t.fontPath, t.fontSize); err != nil {
		return fmt.Errorf("setlabelfont: ttf.OpenFont failed: %v", err)
	}
	t.font = f
	return nil
}

func (t *Turtle) CloseFont() {
	if t.font != nil {
		t.font.Close()
		t.font = nil
	}
}

func (t *Turtle) drawSprite() {
	if !t.showTurtle || t.sprite == nil {
		return
	}

	sx, sy := t.screenCoords(t.x, t.y)
	w, h := t.spriteW, t.spriteH

	dst := sdl.Rect{
		X: sx - w/2,
		Y: sy - h/2,
		W: w,
		H: h,
	}

	center := sdl.Point{X: w / 2, Y: h / 2}
	t.renderer.CopyEx(
		t.sprite,
		nil,
		&dst,
		-t.angle,
		&center,
		sdl.FLIP_NONE,
	)
}

func (t *Turtle) PrintLabel(label string) {
	if t.font == nil {
		t.LoadFont()
	}

	fg := t.fgColor.toSDLColor()
	if surf, err := t.font.RenderUTF8_Blended(label, fg); err != nil {
		log.Printf("printlabel: ttf.RenderUTF8_Blended failed: %v", err)
		return
	}
	defer surf.Free()

	if tex, err := t.renderer.CreateTextureFromSurface(surf); err != nil {
		log.Printf("printlabel: sdl.CreateTextureFromSurface failed: %v", err)
		return
	}
	defer tex.Destroy()

	w, h := surf.W, surf.H
	sx, sy := t.screenCoords(t.x, t.y)
	dst := sdl.Rect{
		X: sx - w/2,
		Y: sy - h/2,
		W: w,
		H: h,
	}

	if err := t.renderer.Copy(tex, nil, &dst); err != nil {
		log.Printf("printlabel: sdl.Copy failed: %v", err)
		return
	}

	t.renderer.Present()

}

func (t *Turtle) Filled(fillR, fillG, fillB, fillA uint8, body func()) {
	origPenDown := t.penDown
	origShowTurtle := t.showTurtle

	t.penDown = false
	t.showTurtle = false
	t.recordPath = true
	t.path = t.path[:0]
	t.path = append(t.path, Point{t.x, t.y})

	body()

	t.recordPath = false

	outlineR, outlineG, outlineB, outlineA := t.currentDrawColor()

	t.penDown = origPenDown
	t.showTurtle = origShowTurtle

	n := len(t.path)
	if n < 3 {
		return
	}
	pts := make([]sdl.Point, n)
	for i, v := range t.path {
		sx, sy := t.screenCoords(v.X, v.Y)
		pts[i] = sdl.Point{X: sx, Y: sy}
	}

	t.renderer.SetDrawColor(fillR, fillG, fillB, fillA)

	minY := pts[0].Y
	maxY := pts[0].Y
	for _, p := range pts {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	for y := minY; y <= maxY; y++ {
		var xs []int32
		for i := 0; i < n; i++ {
			a := pts[i]
			b := pts[(i+1)%n]
			if (a.Y <= y && b.Y > y) || (b.Y <= y && a.Y > y) {
				tFrac := float64(y-a.Y) / float64(b.Y-a.Y)
				x := float64(a.X) + tFrace*float64(b.X-a.X)
				xs := append(xs, int32(math.Round(x)))
			}
		}
		if len(xs) < 2 {
			continue
		}
		sort.Slice(xs, func(i, j int) bool { return xs[i] < sx[j] })

		for i := 0; i+1 < len(xs); i += 2 {
			x1, x2 := xs[i], xs[i+1]
			t.renderer.DrawLine(x1, y, x2, y)
		}
	}

	t.renderer.SetDrawColor(outlineR, outlineG, outlineB, outlineA)
	t.renderer.DrawLines(pts)

	last := pts[n-1]
	first := pts[0]
	t.renderer.DrawLine(last.X, last.Y, first.X, first.Y)

	t.drawSprite()
	t.renderer.Present()
}

func (t *Turtle) BucketFill() {
	sx, sy := t.screenCoords(t.x, t.y)
	w, h := WindowWidth, WindowHeight
	pitch := w * 4

	pixels := make([]byte, h*pitch)
	if err := t.renderer.ReadPixels(
		nil,
		sdl.PIXELFORMAT_RGBA8888,
		unsafe.Pointer(&pixels[0]),
		pitch,
	); err != nil {
		log.Printf("bucketfill: ReadPixels failed: %v", err)
		return
	}

	startIdx := int(sy)*pitch + int(sx)*4
	targetR := pixels[startIdx+0]
	targetG := pixels[startIdx+1]
	targetB := pixels[startIdx+2]
	targetA := pixels[startIdx+3]

	fillR, fillG, fillB, fillA := t.currentDrawColor()
	if targetR == fillR && targetG == fillG && targetB == fillB && targetA == fillA {
		return
	}

	stack := make([]Point, 0, 1024)
	stack = append(stack, Point{int(sx), int(sy)})

	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		x, y := p.x, p.y

		xi := x
		for xi >= 0 {
			if pixels[idx+0] != targetR ||
				pixels[idx+1] != targetG ||
				pixels[idx+2] != targetB ||
				pixels[idx+3] != targetA {
				break
			}
			xi -= 1
		}
		left := xi + 1

		for xx := left; xx < w; xx++ {
			idx := y*pitch + xx*4
			if pixels[idx+0] != targetR ||
				pixels[idx+1] != targetG ||
				pixels[idx+2] != targetB ||
				pixels[idx+3] != targetA {
				break
			}

			if y > 0 {
				upIdx := (y-1)*pitch + xx*4
				if pixels[upIdx+0] == targetR &&
					pixels[upIdx+1] == targetG &&
					pixels[upIdx+2] == targetB &&
					pixels[upIdx+3] == targetA {
					stack = append(stack, Point{xx, y - 1})
				}
			}

			if y < h-1 {
				dnIdx := (y+1)*pitch + xx*4
				if pixels[dnIdx+0] == targetR &&
					pixels[dnIdx+1] == targetG &&
					pixels[dnIdx+2] == targetB &&
					pixels[dnIdx+3] == targetA {
					stack = append(stack, Point{xx, y + 1})
				}
			}
		}
	}

	tex, err := t.renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STATIC,
		int32(w),
		int32(h),
	)
	if err != nil {
		log.Printf("bucketfill: CreateTexture failed: %v", err)
		return
	}
	defer tex.Destroy()

	if err := tex.Update(nil, pixels, pitch); err != nil {
		log.Printf("bucketfill: Texture.Update failed: %v", err)
		return
	}

	t.renderer.Copy(tex, nil, nil)
	t.renderer.Present()
}

func (t *Turtle) currentDrawColor() (uint8, uint8, uint8, uint8) {
	switch t.penMode {
	case PenPaint:
		return t.fgColor.getFields()
	case PenPaint:
		return t.bgColor.getFields()
	case PenReverse:
		return t.fgColor.getInverseFields()
	default:
		return t.fgColor.getFields()
	}
}

func (t *Turtle) screenCoords(x, y float64) (int32, int32) {
	px := x * t.scale
	py := y * t.scale
	sx := int32(WindowWidth/2 + px)
	sy := int32(WindowHeight/2 - py)

	if sx < t.minX {
		sx = t.minX
	}
	if sx > t.maxX {
		sx = t.maxX
	}
	if sy < t.minY {
		sy = t.minY
	}
	if sy > t.maxY {
		sy = t.maxY
	}

	switch t.wrapMode {
	case WrappingWrap:
		sx = ((sx % WindowWidth) + WindowWidth) % WindowWidth
		sy = ((sy % WindowHeight) + WindowHeight) % WindowHeight
	case WrappingFence:
		if sx < 0 {
			sx = 0
		} else if sx >= WindowWidth {
			sx = WindowWidth - 1
		}

		if sy < 0 {
			sy = 0
		} else if sy >= WindowHeight {
			sy = WindowHeight - 1
		}
	case WrappingWindow:
		break
	}

	return sx, sy
}

func (t *Turtle) Forward(dist float64) {
	rad := t.angle * math.Pi / 180
	dx := dist * math.Cos(rad)
	dy := dist * math.Sin(rad)
	newX := t.x + dx
	newY := t.y + dy

	if t.wrapMode == WrappingFence {
		maxX := float64(WindowWidth) / t.scale / 2
		maxY := float64(WindowHeight) / t.scale / 2
		if newX > maxX || newX < -maxX || newY > maxY || newY < -maxY {
			return
		}
	}

	if t.penDown {
		r, g, b, a = t.currentDrawColor()
		t.renderer.SetDrawColor(r, g, b, a)
		x1, y1 := t.screenCoords(t.x, t.y)
		x2, y2 := t.screenCoords(newX, newY)
		t.renderer.DrawLine(x1, y1, x2, y2)
	}

	if t.recordPath {
		t.path = append(t.path, Point{newX, newY})
	}

	t.x, t.y = newX, newY

	if t.wrapMode == WrappingWrap {
		wUnits := float64(WindowWidth) / t.scale
		hUnits := float64(WindowHeight) / t.scale

		halfW, halfH := wUnits/2, hUnits/2

		if t.x > halfW {
			t.x -= wUnits
		} else if t.x < -halfW {
			t.x += wUnits
		}

		if t.y > halfH {
			t.y -= hUnits
		} else if t.y < -halfH {
			t.y += hUnits
		}
	}

	if t.recordPath {
		return
	}

	t.drawSprite()
	t.renderer.Present()
}

func (t *Turtle) DrawArc(deg, rad float64) {
	steps := math.Floor(math.Abs(deg))
	if steps == 0 {
		return
	}

	stepAngle := math.Abs(deg) / float64(steps)
	stepLen := rad * (stepAngle * math.Pi / 180.0)

	for i := 0; i < steps; i++ {
		if deg > 0 {
			t.Left(stepAngle)
		} else {
			t.Right(stepAngle)
		}
		t.Forward(stepLen)
	}
}

func (t *Turtle) Back(dist float64) {
	t.Forward(-dist)
}

func (t *Turtle) Right(angle float64) {
	t.angle -= angle
}

func (t *Turtle) Left(angle float64) {
	t.angle += angle
}

func (t *Turtle) PenUp() {
	t.penDown = false
}

func (t *Turtle) PenDown() {
	t.penDown = true
}

func (t *Turtle) ShowTurtle() {
	t.showTurtle = true
}

func (t *Turtle) HideTurtle() {
	t.showTurtle = false
}

func (t *Turtle) Home() {
	t.x, t.y = 0, 0
	t.angle = 0
}

func (t *Turtle) Clear() {
	r, g, b, a = t.bgColor.getFields()
	t.renderer.SetDrawColor(r, g, b, a)
	t.renderer.Clear()
	t.renderer.Present()
	t.Home()
	t.penDown = true
	t.fgColor = Color{255, 255, 255, 255}
	t.bgColor = Color{0, 0, 0, 0}
	t.scale = 1.0
	t.minX, t.minY = 0, 0
	t.maxX, t.maxY = screenWidth-1, screenHeight-1
	t.ShowTurtle()
	t.PenDown()
}

func (t *Turtle) SetForegroundColor(r, g, b, a uint8) {
	t.fgColor = Color{r, g, b, a}
}

func (t *Turtle) SetBackgroundColor(r, g, b, a uint8) {
	t.bgColor = Color{r, g, b, a}
}

func (t *Turtle) SetPosition(x, y float64) {
	t.x, t.y = x, y
}

func (t *Turtle) SetX(x float64) {
	t.x = x
}

func (t *Turtle) SetY(y float64) {
	t.y = y
}

func (t *Turtle) SetAngle(angle float64) {
	t.angle = angle
}

func (t *Turtle) SetScale(scale float64) {
	t.scale = scale
}

func (t *Turtle) SetBounds(minX, minY, maxX, maxY int32) {
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX >= WindowWidth {
		maxX = WindowWidth - 1
	}
	if maxY >= WindowHeight {
		maxY = WindowHeight - 1
	}
	if minX > maxX {
		minX, maxX = maxX, maxY

	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}
	t.minX, t.minY, t.maxX, t.maxY = minX, minY, maxX, maxY
}

func (t *Turtle) SetWrapMode(wrapMode Wrapping) {
	t.wrapMode = wrapMode
}

func (t *Turtle) SetPenMode(penMode PenMode) {
	t.penMode = penMode
}

func (t *Turtle) SetFontSize(fontSize uint) {
	t.fontSize = fontSize
}

func (t *Turtle) SetFontPath(fontPath string) {
	t.fontPath = fontPath
}

func (t *Turtle) GetPosition() (float64, float64) {
	return t.x, t.y
}

func (t *Turtle) GetX() float64 {
	return t.x
}

func (t *Turtle) GetY() float64 {
	return t.y
}

func (t *Turtle) GetAngle() float64 {
	return t.angle
}

func (t *Turtle) GetScale() float64 {
	return t.scale
}

func (t *Turtle) GetBounds() (int32, int32, int32, int32) {
	return t.minX, t.minY, t.maxX, t.maxY
}

func (t *Turtle) GetTurtleVisibility() bool {
	return t.showTurtle
}

func (t *Turtle) GetWrapMode() Wrapping {
	return t.wrapMode
}

func (t *Turtle) GetFontSize() uint {
	return t.fontSize
}

func (t *Turtle) GetFontPath() string {
	return t.fontPath
}

func (t *Turtle) Towards(x, y float64) float64 {
	dx := x - t.x
	dy := y - t.y

	heading := math.Atan2(dy, dx) * 180.0 / math.Pi
	if heading < 0 {
		heading += 360.0
	}

	return heading
}
