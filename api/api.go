package api

import (
	"postcodes/area"
	"postcodes/service/postcodesio"
)

// api represents the http api, and has the necessary data to response.
type api struct {
	areas *area.Areas
	api   postcodesio.API
}
