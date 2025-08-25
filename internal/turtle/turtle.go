package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Wrapping int

const (
	WrappingWrap   Wrapping = iota
	WrappingFence  Wrapping
	WrappingWindow Wrapping
)

var (
	WindowWidth  = 800
	WindowHeight = 600
)

type Turtle struct {
	x, y       float64
	angle      float64
	penDown    bool
	showTurtle bool
	wrapMode   Wrapping
	r, g, b, a uint8
	scale      float64
	minX, minY int32
	maxX, maxY int32
	spriteW    int32
	spriteH    int32
	renderer   *sdl.Renderer
	sprite     *sdl.Texture
}

func NewTurtle(r *sdl.Renderer, s *sdl.Texture) *Turtle {
	t := &Turtle{
		x:          0,
		y:          0,
		angle:      0,
		penDown:    false,
		showTurtle: false,
		wrapMode:   WrappingWrap,
		r:          255,
		g:          255,
		b:          255,
		a:          255,
		scale:      1.0,
		minX:       0,
		minY:       0,
		maxX:       WindowWidth - 1,
		maxY:       WindowHeight - 1,
		spriteW:    -1,
		spriteH:    -1,
		renderer:   r,
		sprite:     s,
	}
	return t
}

func LoadTurtleImage(path string) error {
	tex, err := img.LoadTexture(t.renderer, path)
	if err != nil {
		return err
	}

	var query sdl.TextureInfo
	if err := tex.Query(&query); err != nil {
		tex.Destroy()
		return err
	}

	t.sprite = tex
	t.spriteW = query.Width
	t.spriteH = query.Height

	return nil
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

	fillR, fillG, fillB, fillA := t.r, t.g, t.b, t.a
	if targetR == fillR && targetG == fillG && targetB == fillB && targetA == fillA {
		return
	}

	type point struct{ x, y int }

	stack := make([]point, 0, 1024)
	stack = append(stack, point{int(sx), int(sy)})

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
					stack = append(stack, point{xx, y - 1})
				}
			}

			if y < h-1 {
				dnIdx := (y+1)*pitch + xx*4
				if pixels[dnIdx+0] == targetR &&
					pixels[dnIdx+1] == targetG &&
					pixels[dnIdx+2] == targetB &&
					pixels[dnIdx+3] == targetA {
					stack = append(stack, point{xx, y + 1})
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
		t.renderer.SetDrawColor(t.r, t.g, t.b, t.a)
		x1, y1 := t.screenCoords(t.x, t.y)
		x2, y2 := t.screenCoords(newX, newY)
		t.renderer.DrawLine(x1, y1, x2, y2)
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

func (t *Turtle) SetColor(r, g, b, a uint8) {
	t.r, t.g, t.b, t.a = r, g, b, a
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

func (t *Turtle) Home() {
	t.x, t.y = 0, 0
	t.angle = 0
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

func (t *Turtle) GetWrapMode() Wrapping {
	return t.wrapMode
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

func (t *Turtle) Clear() {
	t.renderer.SetDrawColor(0, 0, 0, 255)
	t.renderer.Clear()
	t.renderer.Present()
	t.Home()
	t.penDown = true
	t.r, t.g, t.b, t.a = 255, 255, 255, 255
	t.scale = 1.0
	t.minX, t.minY = 0, 0
	t.maxX, t.maxY = screenWidth-1, screenHeight-1
}
