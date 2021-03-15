package utils

import (
	"math"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

const Epsilon = 2.220446049250313e-16

const EdgeClipError = 2.25 * Epsilon

const FaceClipError = 9.0 * (1.0 / math.Sqrt2) * Epsilon

const CellPadding = 2.0 * (FaceClipError + EdgeClipError)

// r2PointsFromEdge Given an edge retrieves the vertices as r2 point of the edge.
func R2PointsFromEdge(e s2.Edge) (r2.Point, r2.Point) {
	face := Face(e.V0.Vector)
	a, b := VectorToPoints(face, e.V0.Vector)
	c, d := VectorToPoints(face, e.V1.Vector)
	return r2.Point{X: a, Y: b}, r2.Point{X: c, Y: d}
}

// VectorToPoints Given a face and a vector return the pair (u, v)
func VectorToPoints(face int, r r3.Vector) (float64, float64) {
	switch face {
	case 0:
		return r.Y / r.X, r.Z / r.X
	case 1:
		return -r.X / r.Y, r.Z / r.Y
	case 2:
		return -r.X / r.Z, -r.Y / r.Z
	case 3:
		return r.Z / r.X, r.Y / r.X
	case 4:
		return r.Z / r.Y, -r.X / r.Y
	}
	return -r.Y / r.Z, -r.X / r.Z
}

// Face Determine the face of the cube on the vector(point) resides.
func Face(r r3.Vector) int {
	f := r.LargestComponent()
	switch {
	case f == r3.XAxis && r.X < 0:
		f += 3
	case f == r3.YAxis && r.Y < 0:
		f += 3
	case f == r3.ZAxis && r.Z < 0:
		f += 3
	}
	return int(f)
}
