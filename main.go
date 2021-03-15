package main

import (
  "io/ioutil"
  "log"
  "postcodes/api"
  "postcodes/area"
  "postcodes/service/postcodesio"
)

func main() {
  postcodesioApi := postcodesio.API{HttpClient: postcodesio.Client()}
  path := "area/file"
  files, err := ioutil.ReadDir(path)
  if err != nil {
    log.Fatalf("error reading area files folder: %s", err.Error())
  }
  var totalAreas area.Areas
  for _, file := range files {
    areas := area.FromGeoJSONFile(path + "/" + file.Name()).HydrateFromApi(postcodesioApi)
    totalAreas = append(totalAreas, areas...)
  }
  api.New(totalAreas).Run()
}
