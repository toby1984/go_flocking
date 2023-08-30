package main

import "github.com/toby1984/go_vectors/vector2"

func Advance(w *World) {
	w.Visit(func(b *Boid) {
		params := w.params

		newAcceleration := flock(b, params, w)

		newVelocity := b.velocity.Add(&newAcceleration).Limit(params.MaxSpeed)
		newLocation := b.location.Add(newVelocity).WrapIfNecessary(params.ModelMax)

		b.location = *newLocation
		b.velocity = *newVelocity
		b.acceleration = newAcceleration
	})
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
	pos := b.location
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

	mean.Add(cohesionVec.Normalize().Multiply(parameters.CohesionWeight))
	mean.Add(alignmentVec.Normalize().Multiply(parameters.AlignmentWeight))
	mean.Add(separationVec.Normalize().Multiply(parameters.SeparationWeight))
	mean.Add(borderForce.Multiply(parameters.BorderForceWeight))

	return mean
}

func steerTo(params SimulationParams, boid *Boid, target vector2.Vector2) *vector2.Vector2 {

	desiredDirection := target.Minus(boid.location)
	var distance = float64(desiredDirection.Len())
	if distance > 0 {
		desiredDirection = desiredDirection.Normalize()
		if distance < 100 {
			desiredDirection = desiredDirection.Multiply(params.MaxSpeed * float32(distance/100.0))
		} else {
			desiredDirection = desiredDirection.Multiply(params.MaxSpeed)
		}

		desiredDirection = desiredDirection.Minus(boid.velocity)
		desiredDirection = desiredDirection.Limit(params.MaxSteeringForce)
		return desiredDirection
	}
	return &vector2.Vector2{}
}
