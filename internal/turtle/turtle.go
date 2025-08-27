package turtle

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

type point struct {
	X, Y float64
}

type color struct {
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
	bgColor    color
	fgColor    color
	scale      float64
	minX, minY int32
	maxX, maxY int32
	spriteW    int32
	spriteH    int32
	penSize    int32
	fontSize   uint
	fontPath   string
	path       []point
	renderer   *sdl.Renderer
	sprite     *sdl.Texture
	font       *ttf.Font
}

func (c color) toSDLColor() {
	return sdl.color{R: c.R, G: c.G, B: c.B, A: c.A}
}

func (c color) getFields() (uint8, uint8, uint8, uint8) {
	return c.r, c.g, c.b, c.a
}

func (c color) getInverseFields() (uint8, uint8, uint8) {
	return 255 - c.r, 255 - c.g, 255 - c.b, c.a
}

func (p point) toSDLpoint() {
	return sdl.point{X: p.X, Y: p.Y}
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
		bgColor:    color{255, 255, 255, 255},
		fgColor:    color{0, 0, 0, 0},
		scale:      1.0,
		minX:       0,
		minY:       0,
		maxX:       WindowWidth - 1,
		maxY:       WindowHeight - 1,
		spriteW:    -1,
		spriteH:    -1,
		penSize:    1,
		fontSize:   12,
		fontPath:   os.GetEnv("GORTLE_DEFAULT_FONTPATH"),
		path:       make([]point, 0, 1024),
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

	center := sdl.point{X: w / 2, Y: h / 2}
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
	origRecordPath := t.recordPath

	t.penDown = false
	t.showTurtle = false
	t.recordPath = true
	t.path = t.path[:0]
	t.path = append(t.path, point{t.x, t.y})

	body()

	t.recordPath = origRecordPath
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
	minX, minY, maxX, maxY := t.getPolygonBounds(pts)
	t.fillPolygonScanline(pts, minX, minY, maxX, maxY, fillR, fillG, fillB, fillA)
	t.renderer.SetDrawColor(outlineR, outlineG, outlineB, outlineA)
	if t.penSize <= 1 {
		t.renderer.DrawLines(pts)
	} else {
		t.renderer.DrawThickLines(pts, t.penSize)
	}

	if t.penSize <= 1 {
		t.renderer.DrawLine(pts[n-1].X, pts[n-1].Y, pts[0].X, pts[0].Y)
	} else {
		t.renderer.DrawThickLine(pts[n-1].X, pts[n-1].Y, pts[0].X, pts[0].Y, t.penSize)
	}

	t.drawSprite()
	t.renderer.Present()
}

func (t *Turtle) getPolygonBounds(pts []sdl.Point) (minX, minY, maxX, maxY int32) {
	if len(pts) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY = pts[0].X, pts[0].Y
	maxX, maxY = pts[0].X, pts[0].Y

	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return minX, minY, maxX, maxY
}

func (t *Turtle) fillPolygonScanline(pts []sdl.Point, minX, minY, maxX, maxY int32, r, g, b, a uint8) {
	for y := minY; y <= maxY; y++ {
		intersections := make([]int32, 0, 16)

		n := len(pts)
		for i := 0; i < n; i++ {
			p1 := pts[i]
			p2 := pts[(i+1)%n]

			if (p1.Y <= y && p2.Y > y) || (p2.Y <= y && p1.Y > y) {
				if p1.Y != p2.Y {
					tFrac := float64(y-p1.Y) / float64(p2.Y-p1.Y)
					x := float64(p1.X) + tFrac*float64(p2.X-p1.X)
					intersections = append(intersections, int32(math.Round(x)))
				}
			}
		}

		sort.Slice(intersections, func(i, j int) bool {
			return intersections[i] < intersections[j]
		})

		for i := 0; i+1 < len(intersections); i += 2 {
			x1 := intersections[i]
			x2 := intersections[i+1]

			if x1 < minX {
				x1 = minX
			}
			if x2 > maxX {
				x2 = maxX
			}

			if x1 <= x2 {
				t.renderer.SetDrawColor(r, g, b, a)
				t.renderer.DrawLine(x1, y, x2, y)
			}
		}
	}
}

func (t *Turtle) BucketFill() {
	sx, sy := t.screenCoords(t.x, t.y)

	if sx < 0 || sx >= int32(WindowWidth) || sy < 0 || sy >= int32(WindowHeight) {
		return
	}

	fillR, fillG, fillB, fillA := t.currentDrawColor()

	pitch := WindowWidth * 4
	pixels := make([]byte, WindowHeight*pitch)

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

	targetR := pixels[startIdx]
	targetG := pixels[startIdx+1]
	targetB := pixels[startIdx+2]
	targetA := pixels[startIdx+3]

	if targetR == fillR && targetG == fillG && targetB == fillB && targetA == fillA {
		return
	}

	t.scanlineFloodFill(pixels, int(sx), int(sy),
		targetR, targetG, targetB, targetA,
		fillR, fillG, fillB, fillA,
		pitch)

	tex, err := t.renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STATIC,
		int32(WindowWidth),
		int32(WindowHeight),
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

func (t *Turtle) scanlineFloodFill(pixels []byte, x, y int,
	targetR, targetG, targetB, targetA,
	fillR, fillG, fillB, fillA uint8,
	pitch int) {

	type segment struct {
		y, x1, x2, dy int
	}

	stack := make([]segment, 0, 1024)

	matchesTarget := func(px, py int) bool {
		if px < 0 || px >= WindowWidth || py < 0 || py >= WindowHeight {
			return false
		}
		idx := py*pitch + px*4
		return pixels[idx] == targetR &&
			pixels[idx+1] == targetG &&
			pixels[idx+2] == targetB &&
			pixels[idx+3] == targetA
	}

	setPixel := func(px, py int) {
		if px >= 0 && px < WindowWidth && py >= 0 && py < WindowHeight {
			idx := py*pitch + px*4
			pixels[idx] = fillR
			pixels[idx+1] = fillG
			pixels[idx+2] = fillB
			pixels[idx+3] = fillA
		}
	}

	x1 := x
	for x1 >= 0 && matchesTarget(x1, y) {
		x1--
	}
	x1++

	x2 := x
	for x2 < WindowWidth && matchesTarget(x2, y) {
		x2++
	}
	x2--

	stack = append(stack, segment{y, x1, x2, 1})
	stack = append(stack, segment{y, x1, x2, -1})

	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		yNew := s.y + s.dy
		if yNew < 0 || yNew >= WindowHeight {
			continue
		}

		for x := s.x1; x <= s.x2; x++ {
			if !matchesTarget(x, yNew) {
				continue
			}

			xNew1 := x
			for xNew1 >= 0 && matchesTarget(xNew1, yNew) {
				setPixel(xNew1, yNew)
				xNew1--
			}
			xNew1++

			xNew2 := x + 1
			for xNew2 < WindowWidth && matchesTarget(xNew2, yNew) {
				setPixel(xNew2, yNew)
				xNew2++
			}
			xNew2--

			if s.dy > 0 {
				stack = append(stack, segment{yNew, xNew1, xNew2, 1})
			} else {
				stack = append(stack, segment{yNew, xNew1, xNew2, -1})
			}

			x = xNew2
		}
	}
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
		t.renderer.DrawThickLine(x1, y1, x2, y2, t.penSize)
	}

	if t.recordPath {
		t.path = append(t.path, point{newX, newY})
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
	t.fgColor = color{255, 255, 255, 255}
	t.bgColor = color{0, 0, 0, 0}
	t.scale = 1.0
	t.minX, t.minY = 0, 0
	t.maxX, t.maxY = screenWidth-1, screenHeight-1
	t.ShowTurtle()
	t.PenDown()
}

func (t *Turtle) SetForegroundColor(r, g, b, a uint8) {
	t.fgColor = color{r, g, b, a}
}

func (t *Turtle) SetBackgroundColor(r, g, b, a uint8) {
	t.bgColor = color{r, g, b, a}
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

func (t *Turtle) SetPenSize(penSize uint) {
	t.penSize = penSize
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

func (t *Turtle) GetForegroundColor() (uint8, uint8, uint8, uint8) {
	return t.fgColor.getFields()
}

func (t *Turtle) GetBackgroundColor() (uint8, uint8, uint8, uint8) {
	return t.bgColor.getFields()
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

func (t *Turtle) GetPenSize() uint {
	return t.penSize
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
