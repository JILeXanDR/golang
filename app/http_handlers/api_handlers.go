package http_handlers

import (
	"net/http"
	"github.com/JILeXanDR/golang/external_api"
	"fmt"
	"github.com/JILeXanDR/golang/app/response_writer"
	"github.com/JILeXanDR/golang/app/session"
	"log"
	"math/rand"
	"strconv"
	"encoding/json"
	"crypto/md5"
	"io"
)

func generateSmsCode() int {
	code, err := strconv.Atoi(fmt.Sprintf(
		"%v%v%v%v",
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
	))

	if err != nil {
		return 0 // TODO fix
	}

	return code
}

type requested interface {
	request()
}

type requestConfirmPhone struct {
	Phone string `json:"phone"`
}

func (r requestConfirmPhone) request() {

}

type requestCheckSmsCode struct {
	Code string `json:"code"`
}

// отправка смс для подтверждение номера телефона
// необходимо для создания пользователя и доступа к личному кабинету
func ConfirmPhoneHandler(w http.ResponseWriter, r *http.Request) {

	var body = &requestConfirmPhone{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	log.Printf("%+v\n", body)

	sess, err := session.GetSession(r)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}
	smsCode := generateSmsCode()
	sess.Values["code"] = smsCode
	sess.Save(r, w)

	go external_api.SendSms(body.Phone, fmt.Sprintf("Код подтверждения %v", smsCode))

	smsCodeStr := strconv.Itoa(smsCode)

	h := md5.New()
	io.WriteString(h, smsCodeStr)
	md5String := fmt.Sprintf("%x", h.Sum(nil))

	type response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Md5     string `json:"md5"`
	}

	response_writer.JsonResponse(w, response{
		Message: "Код подтверждения был отправлен на ваш телефон. Это действие единоразовое",
		Code:    smsCode,
		Md5:     md5String,
	}, 200)
}

// проверка смс кода
func CheckCodeHandler(w http.ResponseWriter, r *http.Request) {

	var body = &requestCheckSmsCode{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
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

	log.Println(sess.Values)

	intCode, err := strconv.Atoi(body.Code)
	if err != nil {
		response_writer.InternalServerError(w, err)
		return
	}

	if intCode != sess.Values["code"] {
		response_writer.JsonMessageResponse(w, "Неверный код", 400)
		return
	}

	response_writer.JsonMessageResponse(w, "Успех", 200)
}
