package main

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"
)

const (
	maxNoise = 2.3

	neighborDistance = 80.0
	separationDistance = 20.0
	cohesionDistance = 70.0
	wallDistance = 100.0


	// percentage values for parameters
	cohesionAmount = 0.8
	alignmentAmount = 2.3
	separationAmount = 10.5
	maxSpeed = 120.0

)

// in the init just seed the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

type boid struct {
	ID  int
	pos pixel.Vec
	vel pixel.Vec
	color color.RGBA
	size int
	angle float64
}

func (b *boid) draw(target *image.RGBA) {
	x, y := int(b.pos.X), int(b.pos.Y)
	r := image.Rect(x - b.size, y-b.size, x+b.size, y+b.size)
	draw.Draw(target, r, &image.Uniform{C: b.color}, image.ZP, draw.Src)
}

// TODO: Currently we're allocating a lot of memory per loop -- this is bound to be a problem
// with high boid counts..!
func (b *boid) update(dt float64, boids []boid) {
	b.color = colornames.Lawngreen
	// This could be broken down to a single loop instead of nested looping
	// Could also use a spatial map to make this faster
	// First, get neighbouring boids...
	var neighbors []boid
	//foundNeighbors := false
	for i, bo := range boids {
		if b.ID == i {
			continue
		}
		sep := dist(b.pos, bo.pos)
		if sep <= neighborDistance {
			b.color = colornames.Azure
			neighbors = append(neighbors, bo)
			//foundNeighbors = true
		}
	}
	// ... then, compute the three boid rules
	// cohesion
	v1 := b.cohesion(neighbors)

	// separation
	v2 := b.separation(neighbors)

	// alignment
	v3 := b.alignment(neighbors)

	b.vel = b.vel.Add(v1).Add(v2)
	b.vel = addNoise(b.vel).Add(v3)

	b.vel = clampMaxSpeed(b.vel)



	//b.avoidWalls(dt)

	x, y := b.pos.XY()
	x += b.vel.X * dt
	y += b.vel.Y * dt
	fw, fh := float64(winWidth), float64(winHeight)
	if x < 0.0 {
		x = fw
	}
	if x > fw {
		x = 0.0
	}
	if y < 0.0 {
		y = fh
	}
	if y > fh {
		y = 0.0
	}



	b.pos = pixel.V(x, y)
}

func addNoise(vec pixel.Vec) pixel.Vec {
	noiseComponent := pixel.V(0.0, 0.0)
	noiseComponent.X = (rand.Float64() - 0.5) * maxNoise
	noiseComponent.Y = (rand.Float64() - 0.5) * maxNoise
	return vec.Add(noiseComponent)
}

func clampMaxSpeed(vec pixel.Vec) pixel.Vec {
	if math.Abs(vec.Len()) <= maxSpeed {
		return vec
	}

	return vec.Unit().Scaled(maxSpeed)
}

func (b *boid) cohesion(neighbors []boid) pixel.Vec {
	if len(neighbors) == 0 {
		return pixel.V(0.0,0.0)
	}
	center := pixel.V(0.0, 0.0)
	for _, n := range neighbors {
		center = center.Add(n.pos)
		b.color = colornames.Aqua
	}
	// TODO: there's an error here
	center = div(center, float64(len(neighbors)))
	return center.Sub(b.pos).Unit().Scaled(cohesionAmount)
}

func (b *boid) separation(neighbors []boid) pixel.Vec {
	if len(neighbors) == 0 {
		return pixel.V(0.0,0.0)
	}
	target := b.pos
	found := false
	for _, n := range neighbors {
		dist := dist(b.pos, n.pos)
		if dist <= separationDistance {
			dp := n.pos.Sub(b.pos)
			target = target.Sub(dp)
			b.color = colornames.Darkseagreen
			found = true
		}
	}
	if !found {
		return pixel.V(0.0, 0.0)
	}
	return target.Sub(b.pos).Unit().Scaled(separationAmount)
}

func (b *boid) alignment(neighbors []boid) pixel.Vec {
	if len(neighbors) == 0 {
		return pixel.V(0.0, 0.0)
	}
	target := b.vel
	for _, n := range neighbors {
		d := dist(b.pos, n.pos)
		if d < cohesionDistance {
			target = target.Add(n.vel)
		}
	}
	return target.Sub(b.vel).Unit().Scaled(alignmentAmount)
}

func (b *boid) avoidWalls(dt float64) {
	fw, fh := float64(winWidth), float64(winHeight)
	x, y := b.pos.XY()
	x += b.vel.X * dt
	y += b.vel.Y * dt

	if x + wallDistance > fw || x - wallDistance < 0 {
		b.vel.X = b.vel.X * -1
	}
	if y + wallDistance > fh || y - wallDistance < 0 {
		b.vel.Y = b.vel.Y * -1
	}
}

