package simulation

import (
	"github.com/toby1984/go_vectors/vector2"
)

func Advance(currentWorld *World) *World {
	var nextWorld World
	nextWorld.SimulationParams = currentWorld.SimulationParams

	nextWorld.Init(currentWorld.SimulationParams, false)

	currentWorld.Visit(func(b *Boid) {
		params := currentWorld.SimulationParams

		newAcceleration := flock(b, params, currentWorld)
		newVelocity := b.Velocity.Add(newAcceleration).Limit(params.MaxSpeed)
		newLocation := b.Location.Add(newVelocity).WrapIfNecessary(params.ModelMax)

		nextWorld.Add(b.CreateCopyWith(newAcceleration, newLocation, newVelocity))
	})
	return &nextWorld
}

func flock(b *Boid, parameters SimulationParams, world *World) vector2.Vector2 {

	visitor := NewNeighborAggregator(b, parameters.SeparationRadius)
	b.VisitNeighbours(world, parameters.NeighbourRadius, visitor)

	// cohesion
	cohesionVec := steerTo(parameters, b, visitor.getAverageLocation())

	// alignment
	alignmentVec := visitor.GetAverageVelocity()

	// separation
	separationVec := visitor.GetAverageSeparationHeading()

	// border force
	pos := b.Location
	borderForce := vector2.Vector2{}

	if pos.X < parameters.BorderRadius {
		delta := (parameters.BorderRadius - pos.X) / parameters.BorderRadius
		borderForce.X = delta * delta
	} else if pos.X > (parameters.ModelMax - parameters.BorderRadius) {
		delta := (parameters.BorderRadius - (parameters.ModelMax - pos.X)) / parameters.BorderRadius
		borderForce.X = -(delta * delta)
	}

	if pos.Y < parameters.BorderRadius {
		delta := (parameters.BorderRadius - pos.Y) / parameters.BorderRadius
		borderForce.Y = delta * delta
	} else if pos.Y > (parameters.ModelMax - parameters.BorderRadius) {
		delta := (parameters.BorderRadius - (parameters.ModelMax - pos.Y)) / parameters.BorderRadius
		borderForce.Y = -(delta * delta)
	}

	mean := vector2.Vector2{}

	mean = mean.Add(cohesionVec.Normalize().Multiply(parameters.CohesionWeight))
	mean = mean.Add(alignmentVec.Normalize().Multiply(parameters.AlignmentWeight))
	mean = mean.Add(separationVec.Normalize().Multiply(parameters.SeparationWeight))
	mean = mean.Add(borderForce.Multiply(parameters.BorderForceWeight))

	return mean
}

func steerTo(params SimulationParams, boid *Boid, target vector2.Vector2) vector2.Vector2 {

	desiredDirection := target.Minus(boid.Location)
	var distance = float64(desiredDirection.Len())
	if distance > 0 {
		desiredDirection = desiredDirection.Normalize()
		if distance < 100 {
			desiredDirection = desiredDirection.Multiply(params.MaxSpeed * float32(distance/100.0))
		} else {
			desiredDirection = desiredDirection.Multiply(params.MaxSpeed)
		}

		desiredDirection = desiredDirection.Minus(boid.Velocity)
		desiredDirection = desiredDirection.Limit(params.MaxSteeringForce)
		return desiredDirection
	}
	return vector2.Vector2{}
}
