package simulation

type SimulationParams struct {
	PopulationSize int32 // number of boids to simulate
	// model coordinates maximum
	// X/Y coordinates are (0,...maximum[
	ModelMax float32

	MaxSteeringForce  float32
	MaxSpeed          float32
	CohesionWeight    float32
	SeparationWeight  float32
	AlignmentWeight   float32
	BorderForceWeight float32
	SeparationRadius  float32
	NeighbourRadius   float32
	BorderRadius      float32
}

/*
int populationSize, double modelMax, double maxForce, double maxSpeed,
            double cohesionWeight, double separationWeight, double alignmentWeight, double borderForceWeight,
            double separationRadius, double neightbourRadius, double borderRadius
*/
// new SimulationParameters(10000,2000,5,10,0.33,0.4,0.33,1,20,100, 5000 * 0.1 )

func GetDefaultSimulationParams() SimulationParams {
	return SimulationParams{
		PopulationSize:    10,
		ModelMax:          300,
		MaxSteeringForce:  5,
		MaxSpeed:          10,
		CohesionWeight:    0.33,
		SeparationWeight:  0.4,
		AlignmentWeight:   0.33,
		BorderForceWeight: 1,
		SeparationRadius:  20,
		NeighbourRadius:   100,
		BorderRadius:      5000 * 0.1,
	}
}
