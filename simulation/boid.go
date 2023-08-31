package simulation

import (
	"github.com/toby1984/go_vectors/vector2"
	"sync/atomic"
)

var uniqueId atomic.Int32

type Boid struct {
	Id           int32
	Acceleration vector2.Vector2
	Location     vector2.Vector2
	Velocity     vector2.Vector2
}

func newBoid() Boid {
	var result Boid
	result.Id = uniqueId.Add(1)
	return result
}

func (b *Boid) CreateCopyWith(acc vector2.Vector2, loc vector2.Vector2, vel vector2.Vector2) Boid {
	var result Boid
	result.Id = b.Id
	result.Location = loc
	result.Acceleration = acc
	result.Velocity = vel
	return result
}

func (b *Boid) VisitNeighbours(world *World, neighborRadius float32, visitor NeighborAggregator) {
	pos := b.Location
	world.VisitNearestBoids(pos.X, pos.Y, neighborRadius, visitor.Visit)
}

func (b *Boid) String() string {
	return "Boid[ loc: " + b.Location.String() + ", vel: " + b.Velocity.String() + ", acc: " + b.Acceleration.String()
}
