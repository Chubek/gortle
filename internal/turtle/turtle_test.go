package turtle_test

import (
	"testing"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// TestTurtleInitialization tests the basic initialization of the Turtle
func TestTurtleInitialization(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Check initial state
	if turtle.x != 0 {
		t.Errorf("Expected initial x=0, got %f", turtle.x)
	}
	if turtle.y != 0 {
		t.Errorf("Expected initial y=0, got %f", turtle.y)
	}
	if turtle.heading != 0 {
		t.Errorf("Expected initial heading=0, got %f", turtle.heading)
	}
	if turtle.penDown != true {
		t.Errorf("Expected initial penDown=true, got %v", turtle.penDown)
	}
	if turtle.penSize != 1 {
		t.Errorf("Expected initial penSize=1, got %d", turtle.penSize)
	}
}

// TestTurtleMovement tests basic turtle movement functions
func TestTurtleMovement(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Test Forward
	initialX, initialY := turtle.x, turtle.y
	turtle.Forward(100)
	if turtle.x == initialX && turtle.y == initialY {
		t.Error("Forward() did not change turtle position")
	}

	// Test Backward
	initialX, initialY = turtle.x, turtle.y
	turtle.Backward(50)
	if turtle.x == initialX && turtle.y == initialY {
		t.Error("Backward() did not change turtle position")
	}

	// Test TurnRight
	initialHeading := turtle.heading
	turtle.TurnRight(90)
	if turtle.heading == initialHeading {
		t.Error("TurnRight() did not change turtle heading")
	}

	// Test TurnLeft
	initialHeading = turtle.heading
	turtle.TurnLeft(45)
	if turtle.heading == initialHeading {
		t.Error("TurnLeft() did not change turtle heading")
	}
}

// TestTurtlePenOperations tests pen up/down functionality
func TestTurtlePenOperations(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Test PenUp
	turtle.PenUp()
	if turtle.penDown != false {
		t.Errorf("Expected penDown=false after PenUp, got %v", turtle.penDown)
	}

	// Test PenDown
	turtle.PenDown()
	if turtle.penDown != true {
		t.Errorf("Expected penDown=true after PenDown, got %v", turtle.penDown)
	}

	// Test SetPenSize
	turtle.SetPenSize(5)
	if turtle.penSize != 5 {
		t.Errorf("Expected penSize=5, got %d", turtle.penSize)
	}
}

// TestTurtlePositioning tests position and home functionality
func TestTurtlePositioning(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Move turtle to a new position
	turtle.Forward(100)
	turtle.TurnRight(90)
	turtle.Forward(50)

	// Record current position
	currentX, currentY := turtle.x, turtle.y

	// Test SetXY
	turtle.SetXY(10, 20)
	if turtle.x != 10 || turtle.y != 20 {
		t.Errorf("SetXY failed: expected (10,20), got (%f,%f)", turtle.x, turtle.y)
	}

	// Test SetHeading
	turtle.SetHeading(180)
	if turtle.heading != 180 {
		t.Errorf("SetHeading failed: expected 180, got %f", turtle.heading)
	}

	// Test Home
	turtle.Home()
	if turtle.x != 0 || turtle.y != 0 || turtle.heading != 0 {
		t.Errorf("Home failed: expected (0,0,0), got (%f,%f,%f)", turtle.x, turtle.y, turtle.heading)
	}
}

// TestTurtleFilled tests the Filled function
func TestTurtleFilled(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Test Filled with a simple square
	turtle.Filled(255, 0, 0, 255, func() {
		for i := 0; i < 4; i++ {
			turtle.Forward(50)
			turtle.TurnRight(90)
		}
	})

	// Verify that the path was recorded correctly
	if len(turtle.path) < 4 {
		t.Errorf("Filled did not record enough path points: got %d", len(turtle.path))
	}
}

// TestTurtleClearScreen tests the ClearScreen function
func TestTurtleClearScreen(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Draw something first
	turtle.Forward(100)
	turtle.TurnRight(90)
	turtle.Forward(100)

	// Record path length
	pathLength := len(turtle.path)

	// Clear screen
	turtle.ClearScreen()

	// Verify that the path was cleared
	if len(turtle.path) != 0 {
		t.Errorf("ClearScreen did not clear path: expected 0, got %d", len(turtle.path))
	}

	// Verify turtle is at home position
	if turtle.x != 0 || turtle.y != 0 || turtle.heading != 0 {
		t.Errorf("ClearScreen did not reset turtle position: got (%f,%f,%f)", turtle.x, turtle.y, turtle.heading)
	}
}

// TestTurtleColorOperations tests color setting functions
func TestTurtleColorOperations(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Test SetDrawColor
	turtle.SetDrawColor(255, 128, 64, 255)
	r, g, b, a := turtle.currentDrawColor()
	if r != 255 || g != 128 || b != 64 || a != 255 {
		t.Errorf("SetDrawColor failed: expected (255,128,64,255), got (%d,%d,%d,%d)", r, g, b, a)
	}

	// Test SetFillColor
	turtle.SetFillColor(100, 150, 200, 255)
	// Note: We can't directly check fill color without exposing it, but we can test that it doesn't crash
}

// TestTurtleScreenCoords tests coordinate conversion
func TestTurtleScreenCoords(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Test screen coordinate conversion
	// Center of screen should correspond to turtle position (0,0)
	sx, sy := turtle.screenCoords(0, 0)
	expectedX := int32(WindowWidth / 2)
	expectedY := int32(WindowHeight / 2)
	if sx != expectedX || sy != expectedY {
		t.Errorf("screenCoords(0,0) failed: expected (%d,%d), got (%d,%d)", expectedX, expectedY, sx, sy)
	}

	// Test with offset
	sx, sy = turtle.screenCoords(100, 50)
	// These values depend on the scale factor, so we just check they're different from center
	if sx == expectedX && sy == expectedY {
		t.Error("screenCoords with offset did not change coordinates")
	}
}

// TestTurtleBounds tests boundary conditions
func TestTurtleBounds(t *testing.T) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		t.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		t.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		t.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Move turtle far beyond screen bounds
	turtle.Forward(10000)
	turtle.TurnRight(90)
	turtle.Forward(10000)

	// Convert to screen coordinates - should not panic
	sx, sy := turtle.screenCoords(turtle.x, turtle.y)
	
	// Just verify we got some coordinates (they might be clamped)
	_ = sx
	_ = sy
}

// BenchmarkTurtleForward benchmarks the Forward function
func BenchmarkTurtleForward(b *testing.B) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		b.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		b.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		b.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		turtle.Forward(10)
	}
}

// BenchmarkTurtleFilled benchmarks the Filled function
func BenchmarkTurtleFilled(b *testing.B) {
	// Initialize SDL for testing
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		b.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		b.Fatalf("Failed to create window: %v", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		b.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		turtle.Filled(255, 0, 0, 255, func() {
			for j := 0; j < 4; j++ {
				turtle.Forward(50)
				turtle.TurnRight(90)
			}
		})
	}
}

// Example test for documentation
func ExampleTurtle_Forward() {
	// Initialize SDL for testing
	sdl.Init(sdl.INIT_VIDEO)
	defer sdl.Quit()

	// Create a window and renderer for testing
	window, _ := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	// Create turtle instance
	turtle := NewTurtle(renderer)

	// Move the turtle forward
	turtle.Forward(100)
	
	// Output: Turtle moved forward by 100 units
}

