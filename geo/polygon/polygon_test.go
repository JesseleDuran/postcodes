package polygon

import (
	"os"
	"testing"

	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

var A, B, C, D Polygon

func TestMain(m *testing.M) {
	A = FromJSONFile("testdata/polygonA.json")
	B = FromJSONFile("testdata/polygonB.json")
	C = FromJSONFile("testdata/polygonC.json")
	D = FromJSONFile("testdata/polygonD.json")
	os.Exit(m.Run())
}

// Given two polygons called A and B its coverage cells should be the same.
func TestPolygon_CoverageCell(t *testing.T) {
	if A.CoverageCell().CellID() != D.CoverageCell().CellID() {
		t.Fatalf("Cells should have the same ID")
	}
	if A.CoverageCell().CellID() == B.CoverageCell().CellID() {
		t.Fatalf("Cells should not have the same ID")
	}
}

// TestFromJSONFile tests if a polygon is well done from a json file,
// comparing the vertices of the polygon decoded in cells,
// with file coordinates in cells
func TestFromJSONFile(t *testing.T) {
	coordinates := [][]float64{
		{-64.20822143554688, 10.802281599725214},
		{-64.2041015625, 10.28924740652188},
		{-63.72344970703124, 10.278437569146138},
		{-63.73718261718749, 10.80363055270312},
		{-64.20822143554688, 10.802281599725214},
	}
	for i, point := range B.Decoded.Vertices() {
		ll := s2.LatLngFromDegrees(coordinates[i][1], coordinates[i][0])
		cExpected := s2.CellFromLatLng(ll).ID()
		c := s2.CellFromPoint(point).ID()
		assert.Equal(t, cExpected, c)
	}
	assert.Equal(t, len(coordinates), len(B.Decoded.Vertices()))
}

// TestFromCoordinates tests if a polygon is well done from an array
// of coordinates, comparing the vertices of the polygon decoded in cells,
// with the coordinates in cells.
func TestFromCoordinates(t *testing.T) {
	coordinates := [][]float64{
		{-64.20822143554688, 10.802281599725214},
		{-64.2041015625, 10.28924740652188},
		{-63.72344970703124, 10.278437569146138},
		{-63.73718261718749, 10.80363055270312},
		{-64.20822143554688, 10.802281599725214},
	}
	p := FromCoordinates(coordinates)
	for i, coord := range coordinates {
		ll := s2.LatLngFromDegrees(coord[1], coord[0])
		cExpected := s2.CellFromLatLng(ll).ID()
		c := s2.CellFromPoint(p.Decoded.Vertices()[i]).ID()

		assert.Equal(t, cExpected, c)
	}
	assert.Equal(t, len(coordinates), len(p.Decoded.Vertices()))
}

func TestPolygon_CreateBound(t *testing.T) {
	A.CreateBound()
	projection := s2.NewMercatorProjection(64)
	for _, v := range A.Decoded.Vertices() {
		A.Bound.ContainsPoint(projection.Project(v))
	}
}

// Given a point should determine if the point is inside or outside
func TestPolygon_ContainsPoint(t *testing.T) {
	// an outside point
	out := s2.LatLngFromDegrees(4.676019722108566, -74.0484470129013)
	// an inside point
	ins := s2.LatLngFromDegrees(4.675784473294116, -74.04892176389694)

	if C.ContainsPoint(s2.PointFromLatLng(out)) {
		t.Fatalf("the point %s %s", out.String(), "should be outside of the polygon")
	}
	if !C.ContainsPoint(s2.PointFromLatLng(ins)) {
		t.Fatalf("the point %s %s", ins.String(), "should be inside of the polygon")
	}
}
