package main

import (
	"github.com/toby1984/go_vectors/vector2"
)

type Boid struct {
	acceleration vector2.Vector2
	location     vector2.Vector2
	velocity     vector2.Vector2
}

func (b *Boid) VisitNeighbours(world *World, neighborRadius float32, visitor NeighborAggregator) {
	pos := b.location
	world.VisitNearestBoids(pos.X, pos.Y, neighborRadius, visitor.Visit)
}
