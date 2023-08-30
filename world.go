package main

/*
	private final KDTree<Boid> tree = new KDTree<>();

	// separate list to keep track of all boids that have been added to
	// the kd-tree , required because traversing the tree to collect
	// all boids is too slow
	private final List<Boid> allBoids = new ArrayList<>();

	private final SimulationParameters simulationParameters;

	public World(SimulationParameters simulationParameters) {
*/

import (
	"github.com/toby1984/go_vectors/vector2"
	"math/rand"
	"sync"
)

type World struct {
	params        SimulationParams
	allBoidsMutex sync.Mutex
	allBoids      []Boid
}

type BoidVisitor func(*Boid)

func (w *World) Init() {
	w.params = GetDefaultSimulationParams()

	w.allBoidsMutex.Lock()
	defer w.allBoidsMutex.Unlock()
	for i := int32(0); i < w.params.PopulationSize; i++ {
		w.allBoids = append(w.allBoids, createRandomBoid(w.params))
	}
}

func (w *World) Add(b Boid) {

	w.allBoidsMutex.Lock()
	defer w.allBoidsMutex.Unlock()
	w.allBoids = append(w.allBoids, b)
}

func (w *World) Visit(visitor BoidVisitor) {

	w.allBoidsMutex.Lock()
	defer w.allBoidsMutex.Unlock()
	for _, boid := range w.allBoids {
		tmp := boid
		visitor(&tmp)
	}
}

func (w *World) VisitNearestBoids(x float32, y float32, maxRadius float32, visitor BoidVisitor) {

	w.allBoidsMutex.Lock()
	radiusSqrt := maxRadius * maxRadius
	defer w.allBoidsMutex.Unlock()
	for _, boid := range w.allBoids {
		if boid.location.DistanceToSqrt(x, y) < radiusSqrt {
			tmp := boid
			visitor(&tmp)
		}
	}
}

func createRandomBoid(params SimulationParams) Boid {
	return Boid{acceleration: createRandomPosition(params), location: createRandomAcceleration(params), velocity: createRandomVelocity(params)}
}

func createRandomPosition(params SimulationParams) vector2.Vector2 {
	x := rand.Float32() * params.ModelMax
	y := rand.Float32() * params.ModelMax
	return vector2.Vector2{X: x, Y: y}
}

func createRandomAcceleration(params SimulationParams) vector2.Vector2 {
	x := rand.Float32() * params.MaxSteeringForce
	y := rand.Float32() * params.MaxSteeringForce
	return vector2.Vector2{X: x, Y: y}
}

func createRandomVelocity(params SimulationParams) vector2.Vector2 {
	x := rand.Float32() * params.MaxSpeed
	y := rand.Float32() * params.MaxSpeed
	return vector2.Vector2{X: x, Y: y}
}
