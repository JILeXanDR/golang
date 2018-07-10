package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
)

func FindAddressHandler(w http.ResponseWriter, r *http.Request) {
	var addresses = []string{
		"Добровольського 6",
		"Остафія Дашкевича 3",
		"Паризької Комуни 64",
	}
	common.JsonResponse(w, addresses, 200)
}
