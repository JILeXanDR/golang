package tests

import (
	"testing"
	"github.com/JILeXanDR/golang/external_api"
)

func TestSearchAddress(t *testing.T) {
	places, err := external_api.FindAddresses("Добровольського")
	if err != nil {
		t.Errorf("api error: %v", err.Error())
		return
	}

	if len(places) == 0 {
		t.Errorf("No places found")
	}
}
