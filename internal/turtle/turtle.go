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
		renderer: r,
	}
	return t
}

func (t *Turtle) screenCoords(x, y float64) (int32, int32) {
	sx := int32(windowWidth/2 + x)
	sy := int32(windowHeight/2 - y)
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

func (t *Turtle) Home() {
	t.x, t.y = 0, 0
	t.angle = 0
}

func (t *Turtle) Clear() {
	t.renderer.SetDrawColor(0, 0, 0, 255)
	t.renderer.Clear()
	t.renderer.Present()
	t.Home()
	t.penDown = true
	t.r, t.g, t.b, t.a = 255, 255, 255, 255
}
