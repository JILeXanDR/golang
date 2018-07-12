package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/common"
	"encoding/json"
	"github.com/JILeXanDR/golang/db"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/JILeXanDR/golang/common/sms"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/gorilla/mux"
	"strconv"
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

func getOrder(r *http.Request) (*db.Order, error) {
	vars := mux.Vars(r)

	var (
		id    int
		err   error
		order = &db.Order{}
	)

	id, err = strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}

	err = db.Connection.First(order, id).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

// создание нового заказа клиентом
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var body = parseBody(r)

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
		common.InternalServerError(w, err)
		return
	}

	var text = fmt.Sprintf("Вы успешно создали заказ под номером %v. Ожидайте СМС по дальнейшей обработке заказа.", order.ID)

	go sms.SendSms("380939411685", text)

	common.JsonResponse(w, order, 200)
}

// получение списка всех заказов (менеджером по заказам)
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders = &[]db.Order{}
	err := db.Connection.Find(orders).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}

	common.JsonResponse(w, orders, 200)
}

func handleError(w http.ResponseWriter, err error) {
	if err == gorm.ErrRecordNotFound {
		common.JsonMessageResponse(w, "Заказ не найден", 404)
	} else if err != nil {
		common.InternalServerError(w, err)
	}
}

// отменить заказ (для медежера по заказам)
// отменить можно только заказы в статусе "created"
func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {

	order, err := getOrder(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if order.Status != db.STATUS_CREATED {
		common.JsonMessageResponse(w, "Нельзя изменить статус заказа", 400)
		return
	}

	order.Status = db.STATUS_CANCELED
	order.CancelReason = "отменено менеджером без причины"
	err = db.Connection.Save(order).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}

	go sms.SendSms(order.Phone, fmt.Sprintf("Ваш заказ был отменен. Причина: %v", order.CancelReason))

	common.JsonResponse(w, order, 200)
}

// подтвердить заказ (для медежера по заказам)
// подтвердить можно только заказы со статусом "created"
func ConfirmOrderHandler(w http.ResponseWriter, r *http.Request) {

	order, err := getOrder(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if order.Status != db.STATUS_CREATED {
		common.JsonMessageResponse(w, "Нельзя изменить статус заказа", 400)
		return
	}

	order.Status = db.STATUS_CONFIRMED
	err = db.Connection.Save(order).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}

	go sms.SendSms(order.Phone, "Ваш заказ был подтверждён")

	common.JsonResponse(w, order, 200)
}

// отметить заказ доставленым
// действие применимо только к заказам со статусом "processing"
func DeliverOrderHandler(w http.ResponseWriter, r *http.Request) {

	order, err := getOrder(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if order.Status != db.STATUS_PROCESSING {
		common.JsonMessageResponse(w, fmt.Sprintf("Нельзя изменить статус заказа. Заказ должен быть в статусе '%v', текущий статус '%v'", db.STATUS_PROCESSING, order.Status), 400)
		return
	}

	order.Status = db.STATUS_DELIVERED
	order.DeliveredAt = time.Now()
	err = db.Connection.Save(order).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}

	go sms.SendSms(order.Phone, "Ваш заказ был доставлен")

	common.JsonResponse(w, order, 200)
}

// начать обработку заказа
// действие применимо только к заказам со статусом "confirmed"
func ProcessOrderHandler(w http.ResponseWriter, r *http.Request) {

	order, err := getOrder(r)
	if err != nil {
		handleError(w, err)
		return
	}

	if order.Status != db.STATUS_CONFIRMED {
		common.JsonMessageResponse(w, fmt.Sprintf("Нельзя взять заказ в обработку. Заказ должен быть в статусе '%v', текущий статус '%v'", db.STATUS_CONFIRMED, order.Status), 400)
		return
	}

	order.Status = db.STATUS_PROCESSING
	err = db.Connection.Save(order).Error
	if err != nil {
		common.InternalServerError(w, err)
		return
	}

	go sms.SendSms(order.Phone, "Ваш заказ был взят в обработку")

	common.JsonResponse(w, order, 200)
}
