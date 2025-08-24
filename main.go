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
	WINDOW_Width   = 800
	WINDOW_Height = 600
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
	sx := int32(WINDOW_Width/2 + x)
	sy := int32(WINDOW_Height/2 - y)
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

func interpret(t *Turtle, script []string) {
	for _, line := range script {
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		cmd := strings.ToLower(tokens[0])

		fmt.Println(cmd)

		switch cmd {
		case "forward", "fd":
			d, err := strconv.ParseFloat(tokens[1], 64)
			if err != nil {
				fmt.Println("Bad number: ", tokens[1])
				continue
			}
			t.Forward(d)
		case "back", "bk":
			a, _ := strconv.ParseFloat(tokens[1], 64)
			t.Right(a)
		case "left", "lt":
			a, _ := strconv.ParseFloat(tokens[1], 64)
			t.Left(a)
		case "right":
			a, _ := strconv.ParseFloat(tokens[1], 64)
			t.Right(a)
		case "setcolor":
			r, _ := strconv.Atoi(tokens[1])
			g, _ := strconv.Atoi(tokens[2])
			b, _ := strconv.Atoi(tokens[3])
			t.SetColor(uint8(r), uint8(g), uint8(b), 255)
		case "penup", "pu":
			t.PenUp()
		case "pendown", "pd":
			t.PenDown()
		case "clearscreen":
			t.Clear()
		case "home":
			t.Home()
		default:
			fmt.Println("Unknown command:", cmd)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Gortle",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_Width, WINDOW_Height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("Could not create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Could not create window: %v", err)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()

	turtle := NewTurtle(renderer)

	script := []string{
		"clearscreen",
		"setcolor 255 0 0",
		"pendown",
		"forward 200",
		"right 90",
		"forward 200",
		"right 90",
		"forward 200",
		"right 90",
		"forward 200",
		"right 90",
		"setcolor 0 255 0",
		"penup",
		"forward 100",
		"pendown",
		"left 45",
		"forward 141.4", // diagonal back to center
	}

	interpret(turtle, script)

	for {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE && e.State == sdl.PRESSED {
					return
				}
			}
		}

		sdl.Delay(16)
	}

}
