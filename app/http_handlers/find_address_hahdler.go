package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/JILeXanDR/golang/external_api"
	"github.com/JILeXanDR/golang/app/response_writer"
)

func FindAddressHandler(w http.ResponseWriter, r *http.Request) {

	var query = r.URL.Query().Get("q")
	var addresses = make([]external_api.Address, 0)
	var lastAddresses = make([]string, 0)
	var err error

	if query == "" {
		lastOrders := make([]db.Order, 0)
		// не показывать дубликаты
		// брать телефон с кукис
		err := db.Connection.Find(&lastOrders, &db.Order{Phone: "380939411685"}).Error
		if err != nil {
			response_writer.InternalServerError(w, err)
			return
		}
		for _, order := range lastOrders {
			lastAddresses = append(lastAddresses, order.DeliveryAddress)
			addresses = append(addresses, external_api.Address{
				Id:   order.DeliveryAddressId,
				Name: order.DeliveryAddress,
			})
		}
		// show last addresses
		//addresses = lastAddresses
	} else {
		addresses, err = external_api.FindAddresses(query)
		if err != nil {
			response_writer.HandleError(w, err)
			return
		}
	}

	response_writer.JsonResponse(w, addresses, 200)
}
