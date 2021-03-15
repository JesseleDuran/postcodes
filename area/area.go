package area

import (
	"io/ioutil"
	"log"
	"os"
	"postcodes/geo/polygon"
	"postcodes/service/postcodesio"

	"github.com/golang/geo/s2"
	"github.com/paulmach/go.geojson"
	"github.com/spf13/cast"
)

type Area struct {
	Postcode string
	Name     string
	Polygon  polygon.Polygon
}

type Areas []Area

func FromGeoJSONFile(path string) Areas {
	f, _ := os.Open(path)
	bytes, _ := ioutil.ReadAll(f)
	return FromGeoJSON(bytes)
}

// FromGeoJson create a Area data structure from a geo json value.
func FromGeoJSON(value []byte) Areas {
	fc, _ := geojson.UnmarshalFeatureCollection(value)
	result := make(Areas, 0)
	for _, f := range fc.Features {
		if f.Geometry.IsPolygon() {
			p := polygon.FromCoordinates(f.Geometry.Polygon[0])
			name := cast.ToString(f.Properties["name"])
			result = append(result, Area{Polygon: p, Name: name})
		}
	}
	return result
}

func (area Area) ContainsLatLng(ll s2.LatLng) bool {
	return area.Polygon.ContainsPoint(s2.PointFromLatLng(ll))
}

func (area *Area) SetPostcode(pc string) {
	area.Postcode = pc
}

func (areas Areas) FindByLatLng(ll s2.LatLng) Area {
	for _, a := range areas {
		if a.ContainsLatLng(ll) {
			return a
		}
	}
	return Area{}
}

func (areas Areas) HydrateFromApi(api postcodesio.API) Areas {
	result := make(Areas, 0)
	channel := make(chan Area, len(areas))
	for _, a := range areas {
		area := a
		go func() {
			var postcode string
			var err error
			ll := s2.LatLngFromPoint(area.Polygon.Decoded.Centroid())
			postcode, err = api.PostCode(ll.Lat.String(), ll.Lng.String(), "1")
			if err != nil {
				log.Println("err getting code", err.Error())
			}
			//if there is no postcode in centroid area..
			if postcode == "" {
				postcode = a.HydrateFromApiTesselating(api)
			}
			if postcode == "" {
				log.Printf("error hydrating area: %s, err: %s", area.Name, err)
			}
			area.Postcode = postcode
			channel <- area
		}()
	}
	for range areas {
		a := <-channel
		result = append(result, a)
	}
	return result
}

func (area Area) HydrateFromApiTesselating(api postcodesio.API) string {
	cells := area.Polygon.Tessellate(11)
	batch, channel := make(chan struct{}, 7), make(chan string, len(cells))
	for _, c := range cells {
		batch <- struct{}{}
		cell := c
		go func() {
			ll := cell.LatLng()
			postcode, err := api.PostCode(ll.Lat.String(), ll.Lng.String(), "1")
			log.Println(area.Name, ll.String())
			if err != nil {
				log.Println("err getting code", err.Error())
			}
			channel <- postcode
			<-batch
		}()
	}
	for range cells {
		p := <-channel
		if p != "" {
			return p
		}
	}
	return ""
}
