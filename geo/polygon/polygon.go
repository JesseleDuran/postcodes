// Package polygon expose a set of functions/methods to do operation over
// closed polygons.
package polygon

import (
  "postcodes/utils"

  "github.com/golang/geo/r2"
  "github.com/golang/geo/s2"
)

// Polygon represents a closed loop formed from a set of vertices.
// a closed polygon is one in which the first and the last vertex is the same
// For instance:
// *---*
// |   | -> A rectangle is a closed loop.
// *---*
// Each polygon is composed by:
// 1) Encoded: A encoded representation of a polygon, the encoded value is only
// a string this field is useful for certain operations like, send or receive
// the polygon through a kafka topic.
// 2) Index: An Index for increase the search process of point containment.
// 3) Decoded: Represents the polygon on its decoded value, the decoded value
// is useful for do some operations like get the vertices or create indexes
// for the polygon.
type Polygon struct {
  Encoded string
  Decoded *s2.Loop
  Bound   r2.Rect
  ID      int
}

// FromCoordinates creates a polygon from a set of flat coordinate.
func FromCoordinates(coordinates [][]float64) Polygon {
  points := make([]s2.Point, 0)
  for _, coordinate := range coordinates {
    ll := s2.LatLngFromDegrees(coordinate[1], coordinate[0])
    points = append(points, s2.PointFromLatLng(ll))
  }
  loop := s2.LoopFromPoints(points)
  loop.Normalize()
  polygon := Polygon{Decoded: loop}
  polygon.CreateBound()
  return polygon
}

// CreateBound create's an index for the given polygon.
// You could see the index as a quad-index that starts with the six top-level
// faces of the S2Cell hierarchy and adaptively splits nodes that intersect
// too many edges.
func (polygon *Polygon) CreateBound() {
  bound := r2.EmptyRect()
  for i := 0; i < polygon.Decoded.NumEdges(); i++ {
    e := polygon.Decoded.Edge(i)
    a, b := utils.R2PointsFromEdge(e)
    bound = bound.AddRect(r2.RectFromPoints(a, b))
  }
  polygon.Bound = bound
}

// ContainsPoint Determine if a point is inside of the given polygon.
// The function use ray tracing to determine if the point is outside or
// inside, ray tracing consist on count edge interceptions.
func (polygon Polygon) ContainsPoint(x s2.Point) bool {
  rp := polygon.Decoded.ReferencePoint()
  crosser, count := s2.NewEdgeCrosser(x, rp.Point), 0
  for i := 0; i < polygon.Decoded.NumEdges(); i++ {
    e := polygon.Decoded.Edge(i)
    if crosser.CrossingSign(e.V0, e.V1) != s2.DoNotCross {
      count++
    }
  }
  if rp.Contained {
    return count%2 == 0
  }
  return count%2 != 0
}

// Tessellate represents a polygon as a set of cells.
func (polygon Polygon) Tessellate(level int) []s2.CellID {
  c := s2.CellFromCellID(polygon.CoverageCell().CellID())
  return polygon.tessellate(level, c)
}

func (polygon Polygon) tessellate(level int, cell s2.Cell) []s2.CellID {
  if !polygon.Decoded.IntersectsCell(cell) {
    return []s2.CellID{}
  }
  if cell.Level() >= level {
    return []s2.CellID{cell.ID()}
  }
  result := make([]s2.CellID, 0)
  children, _ := cell.Children()
  for _, child := range children {
    result = append(result, polygon.tessellate(level, child)...)
  }
  return result
}

func (polygon Polygon) CoverageCell() *s2.PaddedCell {
  f, bound := 0, r2.EmptyRect()
  for i := 0; i < polygon.Decoded.NumEdges(); i++ {
    e := polygon.Decoded.Edge(i)
    a, b := utils.R2PointsFromEdge(e)
    f = utils.Face(e.V0.Vector)
    bound = bound.AddRect(r2.RectFromPoints(a, b))
  }
  faceID := s2.CellIDFromFace(f)
  return s2.PaddedCellFromCellID(
    s2.PaddedCellFromCellID(faceID, utils.CellPadding).ShrinkToFit(bound),
    utils.CellPadding,
  )
}
