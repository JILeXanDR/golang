package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"encoding/json"
	"github.com/JILeXanDR/golang/db"
	"github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type deliveryAddress struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type requestOrder struct {
	List            []string        `json:"list"`
	Phone           string          `json:"phone"`
	DeliveryAddress deliveryAddress `json:"delivery_address"`
	Name            string          `json:"name"`
	Comment         string          `json:"comment"`
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

	log.Println(body.DeliveryAddress.Name)

	metadata, err := json.Marshal(body.List)

	var order = &db.Order{
		List:              postgres.Jsonb{metadata},
		DeliveryAddressId: body.DeliveryAddress.Value,
		DeliveryAddress:   body.DeliveryAddress.Name,
		Phone:             body.Phone,
		Name:              body.Name,
		Comment:           body.Comment,
		Status:            db.STATUS_CREATED,
	}
	err = db.Connection.Create(order).Error
	if err != nil {
		common.InternalServerError(w)
		return
	}

	common.JsonResponse(w, order, 200)
}
