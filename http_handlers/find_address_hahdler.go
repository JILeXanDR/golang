package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"github.com/JILeXanDR/golang/maps"
	"github.com/JILeXanDR/golang/db"
)

func FindAddressHandler(w http.ResponseWriter, r *http.Request) {

	var query = r.URL.Query().Get("q")
	var addresses = make([]maps.Address, 0)
	var lastAddresses = make([]string, 0)
	var err error

	if query == "" {
		lastOrders := make([]db.Order, 0)
		// не показывать дубликаты
		// брать телефон с кукис
		err := db.Connection.Find(&lastOrders, &db.Order{Phone: "0939411685"}).Error
		if err != nil {
			common.InternalServerError(w)
			return
		}
		for _, order := range lastOrders {
			lastAddresses = append(lastAddresses, order.DeliveryAddress)
			addresses = append(addresses, maps.Address{
				Id:   order.DeliveryAddressId,
				Name: order.DeliveryAddress,
			})
		}
		// show last addresses
		//addresses = lastAddresses
	} else {
		addresses, err = maps.SearchAddress(query)
		if err != nil {
			common.HandleError(w, err)
			return
		}
	}

	common.JsonResponse(w, addresses, 200)
}
