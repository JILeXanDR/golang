package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/maps"
)

func FindAddressHandler(w http.ResponseWriter, r *http.Request) {

	var query = r.URL.Query().Get("q")
	var places []string
	var err error

	if query == "" {
		// show last addresses
		places = []string{"пока что ничего нет"}
	} else {
		places, err = maps.SearchAddress(query)
		if err != nil {
			common.HandleError(w, err)
			return
		}
	}

	common.JsonResponse(w, places, 200)
}
