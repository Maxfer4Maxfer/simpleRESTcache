package avaisalesru

import (
	"encoding/json"
)

type jsonResp struct {
	Code        string `json:"code"`
	Type        string `json:"type"`
	CountryName string `json:"country_name"`
	CityName    string `json:"city_name"`
	Name        string `json:"name"`
}
type resp struct {
	Slug     string `json:"slug"`
	Subtitle string `json:"subtitle"`
	Title    string `json:"title"`
}

// Parse prepares specified structure for requests from https://places.aviasales.ru/v2/places.json
func Parse(res []byte) ([]byte, error) {

	jsonRes := []jsonResp{}
	json.Unmarshal(res, &jsonRes)

	r := []resp{}

	for _, e := range jsonRes {
		n := resp{
			Slug:  e.Code,
			Title: e.Name,
		}
		switch e.Type {
		case "city":
			n.Subtitle = e.CountryName
		case "airport":
			n.Subtitle = e.CityName
		}
		r = append(r, n)
	}

	return json.Marshal(r)
}
