package api

import (
	"log"
	"net/http"
	"postcodes/coordinate"
	"postcodes/utils"

	"github.com/gin-gonic/gin"
)

// Postcodes is the http handler to return lat,
// lng and postcode from a memory struct.
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
		var p string
		a := (*api.areas).FindByLatLng(c.ToS2LatLng())
		if a.Postcode == "" {
			p, err = api.api.PostCode(utils.FloatToString(c.Lat), utils.FloatToString(c.Lon), "1")
			if err != nil {
				log.Println("error getting postcode:", err.Error())
			}
			a.SetPostcode(p)
		} else {
			p = a.Postcode
		}
		result = append(result, Response{
			Lat:      c.Lat,
			Lng:      c.Lon,
			Postcode: p,
		})
	}
	c.JSON(http.StatusOK, result)
	return
}
