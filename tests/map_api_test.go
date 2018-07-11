package tests

import (
	"testing"
	"github.com/JILeXanDR/golang/maps"
)

func TestSearchAddress(t *testing.T) {
	places, err := maps.SearchAddress("Добровольського")
	if err != nil {
		t.Errorf("api error: %v", err.Error())
		return
	}

	if len(places) == 0 {
		t.Errorf("No places found")
	}
}
