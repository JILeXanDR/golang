package http_handlers

import (
	"net/http"
	"encoding/json"
	"github.com/JILeXanDR/golang/app/db"
	"github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"time"
	"github.com/gorilla/mux"
	"strconv"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/JILeXanDR/golang/external_api"
	"github.com/JILeXanDR/golang/app/response_writer"
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

func handleError(w http.ResponseWriter, err error) {
	if err == gorm.ErrRecordNotFound {
		response_writer.JsonMessageResponse(w, "Заказ не найден", http.StatusNotFound)
	} else if err != nil {
		response_writer.InternalServerError(w, err)
	}
}

// создание нового заказа клиентом
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var body = parseBody(r)
	var phoneExample = "0939411685"

	// TODO do it better
	if len(body.Phone) != len(phoneExample) || !strings.HasPrefix(body.Phone, "0") {
		response_writer.JsonMessageResponse(w, "Номер телефона введен неверно", http.StatusUnprocessableEntity)
		return
	}

	var phone = "38" + body.Phone

	metadata, err := json.Marshal(body.List)

	var order = &db.Order{
		List:              postgres.Jsonb{metadata},
		DeliveryAddressId: body.DeliveryAddress.Value,
		DeliveryAddress:   body.DeliveryAddress.Name,
		Phone:             phone,
		Name:              body.Name,
		Comment:           body.Comment,
		Status:            db.STATUS_CREATED,
	}
	err = db.Connection.Create(order).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	var text = fmt.Sprintf("Вы успешно создали заказ под номером %v. Ожидайте СМС по дальнейшей обработке заказа.", order.ID)

	go external_api.SendSms(order.Phone, text)

	response_writer.JsonResponse(w, order, http.StatusCreated)
}

// получение списка всех заказов (менеджером по заказам)
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orders = &[]db.Order{}
	err := db.Connection.Find(orders).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	response_writer.JsonResponse(w, orders, http.StatusOK)
}

// получение информации о зазазе (для медежера по заказам)
func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	order, err := getOrder(r)
	if err != nil {
		handleError(w, err)
		return
	}

	response_writer.JsonResponse(w, order, http.StatusOK)
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
		response_writer.JsonMessageResponse(w, "Нельзя изменить статус заказа", 400)
		return
	}

	order.Status = db.STATUS_CANCELED
	order.CancelReason = "отменено менеджером без причины"
	err = db.Connection.Save(order).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	go external_api.SendSms(order.Phone, fmt.Sprintf("Ваш заказ был отменен. Причина: %v", order.CancelReason))

	response_writer.JsonResponse(w, order, http.StatusOK)
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
		response_writer.JsonMessageResponse(w, "Нельзя изменить статус заказа", http.StatusBadRequest)
		return
	}

	order.Status = db.STATUS_CONFIRMED
	err = db.Connection.Save(order).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	go external_api.SendSms(order.Phone, "Ваш заказ был подтверждён")

	response_writer.JsonResponse(w, order, http.StatusOK)
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
		response_writer.JsonMessageResponse(w, fmt.Sprintf("Нельзя изменить статус заказа. Заказ должен быть в статусе '%v', текущий статус '%v'", db.STATUS_PROCESSING, order.Status), 400)
		return
	}

	order.Status = db.STATUS_DELIVERED
	order.DeliveredAt = time.Now()
	err = db.Connection.Save(order).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	go external_api.SendSms(order.Phone, "Ваш заказ был доставлен")

	response_writer.JsonResponse(w, order, http.StatusOK)
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
		response_writer.JsonMessageResponse(w, fmt.Sprintf("Нельзя взять заказ в обработку. Заказ должен быть в статусе '%v', текущий статус '%v'", db.STATUS_CONFIRMED, order.Status), 400)
		return
	}

	order.Status = db.STATUS_PROCESSING
	err = db.Connection.Save(order).Error
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	go external_api.SendSms(order.Phone, "Ваш заказ был взят в обработку")

	response_writer.JsonResponse(w, order, http.StatusOK)
}
