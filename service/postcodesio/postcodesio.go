// Package postcodesio is in charge of communication with postcode api.
package postcodesio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"postcodes/service/postcodesio/postcode"
	"time"

	"github.com/go-resty/resty/v2"
)

// API represents a http client.
type API struct {
	HttpClient *resty.Client
}

func Client() API {
	client := resty.New()
	client.SetHeader("X-Application-ID", "postcodes")
	client.SetHostURL(os.Getenv("POSTCODESIO_HOST"))
	client.SetTimeout(10 * time.Second)
	return API{HttpClient: client}
}

// PostCode communicates with postcodesio api to get postcode info given a
// lat, long and a limit of results.
func (api API) PostCode(lat, lon, limit string) (string, error) {
	var pp postcode.Postcodes
	params := map[string]string{
		"lat":   lat,
		"lon":   lon,
		"limit": limit,
	}
	res, err := api.HttpClient.R().SetQueryParams(params).Get(PostcodesURI)
	if err != nil {
		return "", err
	}
	if !res.IsSuccess() {
		return "", errors.New(fmt.Sprintf("invalid response, code %d",
			res.StatusCode()))
	}
	err = json.Unmarshal(res.Body(), &pp)
	if err != nil {
		return "", err
	}
	if len(pp.Result) == 0 {
		return "", nil
	}
	return pp.Result[0].Postcode, nil
}
