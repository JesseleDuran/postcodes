package postcode

//Postcodes is the way that postcodes.io returns postcodes info.
type Postcodes struct {
	Result []struct {
		Postcode string `json:"postcode"`
	} `json:"result"`
}
