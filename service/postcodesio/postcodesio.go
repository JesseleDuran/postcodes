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

type API struct {
  HttpClient *resty.Client
}

func Client() *resty.Client {
  client := resty.New()
  client.SetHeader("X-Application-ID", "postcodes")
  client.SetHostURL(os.Getenv("POSTCODESIO_HOST"))
  client.SetTimeout(10 * time.Second)
  return client
}

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
