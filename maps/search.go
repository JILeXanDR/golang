package maps

import (
	"os"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"github.com/kr/pretty"
)

// возвращает список мест по ключевому слову
func SearchAddress(input string) (places []string, err error) {

	var client *maps.Client
	client, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		return nil, err
	}

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

	for _, prediction := range resp.Predictions {
		pretty.Log(prediction)
		places = append(places, prediction.Description)
	}

	return places, nil
}
