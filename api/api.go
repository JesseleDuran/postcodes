package api

import (
	"postcodes/area"
	"postcodes/service/postcodesio"
)

type api struct {
	areas *area.Areas
	api   postcodesio.API
}
