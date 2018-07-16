package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/external_api"
	"fmt"
	"github.com/JILeXanDR/golang/app/response_writer"
	"github.com/JILeXanDR/golang/app/session"
	"strconv"
	"encoding/json"
	"github.com/JILeXanDR/golang/app/services"
)

type requestJsonHandler interface {
	handle(r *http.Request) error
}

func handle(base interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(base)
	if err != nil {
		return err
	}
	return nil
}

type requestConfirmPhone struct {
	Phone string `json:"phone"`
}

func (base *requestConfirmPhone) handle(r *http.Request) error {
	return handle(base, r)
}

type requestCheckSmsCode struct {
	Code string `json:"code"`
}

func (base *requestCheckSmsCode) handle(r *http.Request) error {
	return handle(base, r)
}

// отправка смс для подтверждение номера телефона
// необходимо для создания пользователя и доступа к личному кабинету
func ConfirmPhoneHandler(w http.ResponseWriter, r *http.Request) {

	var body = &requestConfirmPhone{}
	err := body.handle(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	phone, err := services.NormalizePhone(body.Phone)
	if err != nil {
		response_writer.ValidationError(w, err.Error())
		return
	}

	//log.Printf("%+v\n", body)

	sess, err := session.GetSession(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	smsCode := services.GenerateSmsCode()
	sess.Values["auth"] = services.PhoneSession{Phone: phone, Confirmed: false}
	sess.Values["code"] = smsCode
	sess.Save(r, w)

	go external_api.SendSms(phone, fmt.Sprintf("Код подтверждения %v", smsCode))

	smsCodeMd5 := services.DoubleMd5(strconv.Itoa(smsCode))

	type response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Md5     string `json:"md5"`
	}

	response_writer.JsonResponse(w, response{
		Message: "Код подтверждения был отправлен на ваш телефон. Это действие единоразовое",
		Code:    smsCode,
		Md5:     smsCodeMd5,
	}, 200)
}

// проверка смс кода
func CheckCodeHandler(w http.ResponseWriter, r *http.Request) {

	var body = &requestCheckSmsCode{}
	err := body.handle(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	if len(body.Code) == 0 {
		response_writer.JsonMessageResponse(w, "Код не введен", 422)
		return
	}

	sess, err := session.GetSession(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	intCode, err := strconv.Atoi(body.Code)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	if intCode != sess.Values["code"] {
		response_writer.JsonMessageResponse(w, "Неверный код", 400)
		return
	}

	var phoneSession = sess.Values["auth"].(services.PhoneSession)
	phoneSession.Confirmed = true
	sess.Values["auth"] = phoneSession
	sess.Save(r, w)

	response_writer.JsonMessageResponse(w, "Успех", 200)
}
