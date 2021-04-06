package main

import (
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Snake struct
type Snake struct {
	pos   pixel.Vec
	speed pixel.Vec
	body  []pixel.Vec
}

// Constants
const SNAKE_SPRITE_FILE = "./sprites/guy.png"
const FOOD_SPRITE_FILE = "./sprites/amogus.png"

const HORIZONTAL_PIXELS = 30
const VERTICAL_PIXELS = 20
const FPS = 10
const PIXEL_SIZE = 40

const WINDOW_WIDTH = HORIZONTAL_PIXELS * PIXEL_SIZE
const WINDOW_HEIGHT = VERTICAL_PIXELS * PIXEL_SIZE
const WINDOW_TITLE = "Snake"

// Function to load picture
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	// Create Window
	cfg := pixelgl.WindowConfig{
		Title:  WINDOW_TITLE,
		Bounds: pixel.R(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create snake
	var snake = Snake{window.Bounds().Center(), pixel.V(PIXEL_SIZE, 0), []pixel.Vec{window.Bounds().Center()}}

	// Load snake sprite
	snake_png, err := loadPicture(SNAKE_SPRITE_FILE)
	if err != nil {
		panic(err)
	}

	// Load food sprite
	food_png, err := loadPicture(FOOD_SPRITE_FILE)
	if err != nil {
		panic(err)
	}

	// Convert both to pixel sprites
	snake_sprite := pixel.NewSprite(snake_png, snake_png.Bounds())
	food_sprite := pixel.NewSprite(food_png, food_png.Bounds())

	// Decide position of food (random)
	food_pos := pixel.V(float64(rand.Intn(HORIZONTAL_PIXELS)*PIXEL_SIZE), float64(rand.Intn(VERTICAL_PIXELS)*PIXEL_SIZE))

	// Decide scaling of sprites to fit grid
	snake_sprite_scale := pixel.V(PIXEL_SIZE/snake_png.Bounds().W(), PIXEL_SIZE/snake_png.Bounds().H())
	food_sprite_scale := pixel.V(PIXEL_SIZE/food_png.Bounds().W(), PIXEL_SIZE/food_png.Bounds().H())

	// Get time to decide FPS wait
	last := time.Now()

	// If window isn't closed, i.e. every frame update,
	for !(window.Closed()) {

		// Find time taken to draw one frame
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Handle keyboard events
		if window.Pressed(pixelgl.KeyLeft) {
			snake.speed = pixel.V(-PIXEL_SIZE, 0)
		}
		if window.Pressed(pixelgl.KeyRight) {
			snake.speed = pixel.V(PIXEL_SIZE, 0)
		}
		if window.Pressed(pixelgl.KeyDown) {
			snake.speed = pixel.V(0, -PIXEL_SIZE)
		}
		if window.Pressed(pixelgl.KeyUp) {
			snake.speed = pixel.V(0, PIXEL_SIZE)
		}

		// Sleep to adjust FPS of game
		time.Sleep(time.Duration(float64(time.Second) * (1.0/FPS - dt)))

		// Clear window
		window.Clear(colornames.Black)

		// Add speed vectors to snake and make sure it doesn't escape grid
		snake.pos = snake.pos.Add(snake.speed)
		snake.pos = pixel.V(math.Round(snake.pos.X/PIXEL_SIZE)*PIXEL_SIZE, math.Round(snake.pos.Y/PIXEL_SIZE)*PIXEL_SIZE)

		// Handle reach arounds
		if snake.pos.X > (WINDOW_WIDTH - PIXEL_SIZE) {
			snake.pos.X = 0
		}
		if snake.pos.X < 0 {
			snake.pos.X = WINDOW_WIDTH - PIXEL_SIZE
		}
		if snake.pos.Y > (WINDOW_HEIGHT - PIXEL_SIZE) {
			snake.pos.Y = 0
		}
		if snake.pos.Y < 0 {
			snake.pos.Y = WINDOW_HEIGHT - PIXEL_SIZE
		}

		// Draw food
		food_mat := pixel.IM
		food_mat = food_mat.ScaledXY(pixel.ZV, food_sprite_scale)
		food_mat = food_mat.Moved(food_pos.Add(pixel.V(PIXEL_SIZE/2, PIXEL_SIZE/2)))
		food_sprite.Draw(window, food_mat)

		// Move every part of snake to the part in front of it
		for i := 0; i < (len(snake.body) - 1); i++ {
			snake.body[i] = snake.body[i+1]
		}
		snake.body[len(snake.body)-1] = snake.pos

		// Check if snake ran into itself
		collided := false
		for i := 0; i < (len(snake.body) - 1); i++ {
			if snake.body[i] == snake.pos {
				collided = true
			}
		}

		// If collided, revert snake
		if collided {
			snake.body = []pixel.Vec{snake.pos}
		}

		// Draw every part of snake's body
		for _, cell := range snake.body {
			mat := pixel.IM
			mat = mat.ScaledXY(pixel.ZV, snake_sprite_scale)
			mat = mat.Moved(cell.Add(pixel.V(PIXEL_SIZE/2, PIXEL_SIZE/2)))
			snake_sprite.Draw(window, mat)
		}

		// If snake eats food, add it to its body
		if snake.pos == food_pos {
			food_pos = pixel.V(float64(rand.Intn(HORIZONTAL_PIXELS)*PIXEL_SIZE), float64(rand.Intn(VERTICAL_PIXELS)*PIXEL_SIZE))
			snake.body = append(snake.body, food_pos)
		}

		// Update window
		window.Update()
	}
}

func main() {
	// Run Game from main func
	pixelgl.Run(run)
}
