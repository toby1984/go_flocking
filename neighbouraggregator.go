package main

import (
	"github.com/toby1984/go_vectors/vector2"
	"math"
)

type NeighborAggregator struct {
	separationRadius         float32
	boid                     *Boid
	locationSumX             float32
	locationSumY             float32
	velocitySumX             float32
	velocitySumY             float32
	separationSumX           float32
	separationSumY           float32
	neighbourCount           int32
	separationNeighbourCount int32
}

func NewNeighborAggregator(b *Boid, separationRadius float32) NeighborAggregator {
	var result NeighborAggregator
	result.separationRadius = separationRadius
	result.boid = b
	return result
}

func (v *NeighborAggregator) Visit(otherBoid *Boid) {

	/*
	   TODO: Not the same as in the Java code as there we have a pointer that is being checked for equality
	   if ( boid == otherBoid ) {
	       return;
	   }
	*/
	boid := v.boid

	if v.boid == otherBoid {
		return
	}

	distance := otherBoid.location.DistanceTo(&boid.location)
	v.neighbourCount++

	v.locationSumX += otherBoid.location.X
	v.locationSumY += otherBoid.location.Y

	v.velocitySumX += otherBoid.velocity.X
	v.velocitySumY += otherBoid.velocity.Y

	if distance > 0 && distance < v.separationRadius {
		tmpX := boid.location.X
		tmpY := boid.location.Y

		tmpX -= otherBoid.location.X
		tmpY -= otherBoid.location.Y

		length := tmpX*tmpX + tmpY*tmpY
		if length > 0.00001 {
			length = float32(math.Sqrt(float64(length)))
			tmpX /= length
			tmpY /= length
		}

		v.separationSumX += tmpX
		v.separationSumY += tmpY

		v.separationNeighbourCount++
	}
}

func (v *NeighborAggregator) GetNeighbourCount() int32 {
	return v.neighbourCount
}

// separation
func (v *NeighborAggregator) GetAverageSeparationHeading() vector2.Vector2 {
	if v.separationNeighbourCount == 0 {
		return vector2.Vector2{}
	}
	return vector2.Vector2{X: v.separationSumX / float32(v.separationNeighbourCount),
		Y: v.separationSumY / float32(v.separationNeighbourCount)}
}

// alignment
func (v *NeighborAggregator) GetAverageVelocity() vector2.Vector2 {
	if v.neighbourCount == 0 {
		return vector2.Vector2{}
	}
	return vector2.Vector2{X: v.velocitySumX / float32(v.neighbourCount), Y: v.velocitySumY / float32(v.neighbourCount)}
}

// cohesion
func (v *NeighborAggregator) getAverageLocation() vector2.Vector2 {
	if v.neighbourCount == 0 {
		return vector2.Vector2{}
	}
	return vector2.Vector2{X: v.locationSumX / float32(v.neighbourCount), Y: v.locationSumY / float32(v.neighbourCount)}
}
