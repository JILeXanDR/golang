package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"encoding/json"
	"github.com/JILeXanDR/golang/db"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type requestOrder struct {
	List            []string `json:"list"`
	Phone           string   `json:"phone"`
	DeliveryAddress string   `json:"delivery_address"`
}

func parseBody(r *http.Request) (requestOrder) {
	decoder := json.NewDecoder(r.Body)
	var data requestOrder
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	return data
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {

	var body = parseBody(r)

	metadata, err := json.Marshal(body.List)

	var order = &db.Order{
		List:            postgres.Jsonb{metadata},
		DeliveryAddress: body.DeliveryAddress,
		Phone:           body.Phone,
	}
	err = db.Connection.Create(order).Error
	if err != nil {
		common.InternalServerError(w)
		return
	}

	common.JsonResponse(w, body, 200)
}
