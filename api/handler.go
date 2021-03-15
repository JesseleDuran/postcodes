package api

import (
  "net/http"
  "postcodes/coordinate"

  "github.com/gin-gonic/gin"
)

func (api api) Postcodes(c *gin.Context) {
  req := struct {
    Coordinates coordinate.Coordinates `json:"coordinates"`
  }{}
  err := c.BindJSON(&req)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  type Response struct {
    Lat      float64 `json:"lat"`
    Lng      float64 `json:"lon"`
    Postcode string  `json:"postcode"`
  }
  result := make([]Response, 0)
  for _, c := range req.Coordinates {
    a := api.areas.FindByLatLng(c.ToS2LatLng())
    result = append(result, Response{
      Lat:      c.Lat,
      Lng:      c.Lon,
      Postcode: a.Postcode,
    })
  }
  c.JSON(http.StatusOK, result)
  return
}
