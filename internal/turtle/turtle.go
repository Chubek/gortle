package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
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
	wrap       true
	r, g, b, a uint8
	scale      float64
	minx, miny int32
	maxx, maxy int32
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
		wrap:       true,
		r:          255,
		g:          255,
		b:          255,
		a:          255,
		scale:      1.0,
		minx:       0,
		miny:       0,
		maxx:       WindowWidth - 1,
		maxy:       WindowHeight - 1,
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

func (t *Turtle) screenCoords(x, y float64) (int32, int32) {
	px := x * t.scale
	py := y * t.scale
	sx := int32(WindowWidth/2 + px)
	sy := int32(WindowHeight/2 - py)

	if sx < t.minx {
		sx = t.minx
	}
	if sx > t.maxx {
		sx = t.maxx
	}
	if sy < t.miny {
		sy = t.miny
	}
	if sy > t.maxy {
		sy = t.maxy
	}

	if t.wrap {
		sx = ((sx % WindowWidth) + WindowWidth) % WindowWidth
		sy = ((sy % WindowHeight) + WindowHeight) % WindowHeight
	}

	return sx, sy
}

func (t *Turtle) Forward(dist float64) {
	rad := t.angle * math.Pi / 180
	dx := dist * math.Cos(rad)
	dy := dist * math.Sin(rad)
	newX := t.x + dx
	newY := t.y + dy

	if t.penDown {
		t.renderer.SetDrawColor(t.r, t.g, t.b, t.a)
		x1, y1 := t.screenCoords(t.x, t.y)
		x2, y2 := t.screenCoords(newX, newY)
		t.renderer.DrawLine(x1, y1, x2, y2)
	}

	t.x = newX
	t.y = newY

	if t.wrap {
		wUnits := float64(WindowWidth) / t.scale
		hUnits := float64(WindowHeight) / t.scale

		halfW := wUnits / 2
		halfH := hUnits / 2

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

func (t *Turtle) WrapOn() {
	t.wrap = true
}

func (t *Turtle) WrapOff() {
	t.wrap = false
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

func (t *Turtle) SetBounds(minx, miny, maxx, maxy int32) {
	if minx < 0 {
		minx = 0
	}
	if miny < 0 {
		miny = 0
	}
	if maxx >= WindowWidth {
		maxx = WindowWidth - 1
	}
	if maxy >= WindowHeight {
		maxy = WindowHeight - 1
	}
	if minx > maxx {
		minx, maxx = maxx, maxy
	}
	if miny > maxy {
		miny, maxy = maxy, miny
	}
	t.minx, t.miny, t.maxx, t.maxy = minx, miny, maxx, maxy
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
	return t.minx, t.miny, t.maxx, t.maxy
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
	t.minx, t.miny = 0, 0
	t.maxx, t.maxy = screenWidth-1, screenHeight-1
}
