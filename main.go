package main

import (
	"github.com/toby1984/go_flocking/simulation"
	"github.com/toby1984/go_flocking/tty"
	"time"
)

var term = tty.NewTerminal()

func main() {

	var timerChannel = make(chan int)

	renderFunction := func() {

		frameNo := 0

		var world = new(simulation.World)
		world.Init(simulation.GetDefaultSimulationParams(), true)

		for range timerChannel {
			term.ClearScreen()
			frameNo++
			var cols, rows = term.GetTerminalSize()
			// term.PrintAtWithColor(fmt.Sprintf("Rendering frame #%d (%dx%d)\n", frameNo, cols, rows), 0, 0, tty.COLOR_RED)
			world.Visit(func(b *simulation.Boid) {
				xPerc := b.Location.X / world.SimulationParams.ModelMax
				yPerc := b.Location.Y / world.SimulationParams.ModelMax

				col := xPerc * float32(cols)
				row := yPerc * float32(rows)
				term.PrintAtWithColor("*", int(col), int(row), tty.COLOR_RED)
			})
			world = simulation.Advance(world)
		}
	}
	go renderFunction()

	// timer loop
	for {
		timerChannel <- 1
		time.Sleep(100 * time.Millisecond)
	}
}
