package external_api

import (
	"os"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type Address struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// возвращает список мест по ключевому слову
func FindAddresses(input string) (addresses []Address, err error) {

	var client *maps.Client
	client, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return nil, err
	}

	// координаты центра Черкасс
	var location = &maps.LatLng{
		Lat: 49.421022,
		Lng: 32.056461,
	}

	r := &maps.QueryAutocompleteRequest{
		Input:    input,
		Language: "uk",
		Radius:   10000,
		Location: location,
		Offset:   0,
	}

	resp, err := client.QueryAutocomplete(context.Background(), r)
	if err != nil {
		return nil, err
	}

	addresses = make([]Address, 0)

	for _, prediction := range resp.Predictions {
		addresses = append(addresses, Address{
			Id:   prediction.PlaceID,
			Name: prediction.Description,
		})
	}

	return addresses, nil
}
