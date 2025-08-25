package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

type Turtle struct {
	x, y       float64
	angle      float64
	penDown    bool
	r, g, b, a uint8
	scale      float64
	minx, miny int32
	maxx, maxy int32
	renderer   *sdl.Renderer
}

func NewTurtle(r *sdl.Renderer) *Turtle {
	t := &Turtle{
		x:        0,
		y:        0,
		angle:    0,
		r:        255,
		g:        255,
		b:        255,
		a:        255,
		scale:    1,
		minx:     0,
		miny:     0,
		maxx:     windowWidth - 1,
		maxy:     windowHeight - 1,
		renderer: r,
	}
	return t
}

func (t *Turtle) screenCoords(x, y float64) (int32, int32) {
	px := x * t.scale
	py := y * t.scale
	sx := int32(windowWidth/2 + px)
	sy := int32(windowHeight/2 - py)

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
		t.renderer.Present()
	}

	t.x = newX
	t.y = newY
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
	if maxx >= windowWidth {
		maxx = windowWidth - 1
	}
	if maxy >= windowHeight {
		maxy = windowHeight - 1
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
}
