package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"
	"golang.org/x/image/colornames"
	"image"
	"math"
	"math/rand"
	"time"
)

const (
	winWidth = 1024
	winHeight = 768
)
/// STATS - 1002 BOIDS --> 50-70 FPS ( 2019-09-05 ) - Screen resolution 1024 : 768
///


func main() {
	//TODO: flags for profiling during runtime.
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Destructible Pixel-based Terrain Tests",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())


	boids := createBoids()

	last := time.Now()
	second := time.Tick(time.Second)
	frames := 0
	for !win.Closed() {
		// spawn boids if mouse pressed
		if win.Pressed(pixelgl.MouseButtonLeft) {
			mouse := win.MousePosition()
			b := boid{
				ID: boids[len(boids)-1].ID + 1,
				pos: pixel.V(mouse.X, mouse.Y),
				color: colornames.Lawngreen,
				size: 5,
				angle: rand.Float64() * math.Pi * 180.0,
				vel: pixel.V((rand.Float64()-0.5)*400.0, (rand.Float64()-0.5)*400.0),
			}
			boids = append(boids, b)
		}

		// update delta time
		dt := time.Since(last).Seconds()
		last = time.Now()

		// update objects
		for i := range boids {
			boids[i].update(dt, boids)
		}

		// update buffers
		refreshCanvas(canvas, boids)

		// draw visible stuff
		win.Clear(colornames.Cadetblue)
		canvas.Draw(win,  pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		// Count fps
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | Boids: %d", cfg.Title, frames, len(boids)))
			frames = 0
		default:
		}
	}
}

func createBoids() []boid {
	var boids []boid
	for i := 0; i < 10; i++ {
		x, y := rand.Float64() * winWidth, rand.Float64() *  winHeight
		b := boid{
			ID: i,
			pos: pixel.V(x, y),
			color: colornames.Lawngreen,
			size: 5,
			angle: rand.Float64() * math.Pi * 180.0,
			vel: pixel.V((rand.Float64()-0.5)*400.0, (rand.Float64()-0.5)*400.0),
		}
		boids = append(boids, b)
	}
	return boids
}

func refreshCanvas(canvas *pixelgl.Canvas, boids []boid) {
	buffer := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	for _, b := range boids {
		b.draw(buffer)
	}
	canvas.SetPixels(buffer.Pix)
}
