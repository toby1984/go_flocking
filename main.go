package main

import (
	"github.com/toby1984/go_vectors/vector2"
	"github.com/toby1984/go_vectors/vector3"
	"os"
)

func main() {
	var v1 = vector2.Vector2{X: 1, Y: 2}
	var v2 = vector3.Vector3{Vector2: vector2.Vector2{1, 2}, Z: 3}
	os.Stdout.Write([]byte(v1.ToString()))
	os.Stdout.Write([]byte(v2.ToString()))
}
