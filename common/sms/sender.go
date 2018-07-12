package sms

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"os"
	"encoding/json"
	"net/url"
	"log"
	"time"
)

type apiResponse struct {
	Code    int
	Message string
	data    interface{}
}

type apiRequest struct {
	ApiKey string `json:"api_key"`
}

func GetOwnBalance() {
	uri := fmt.Sprintf("https://api.mobizon.ua/service/user/getownbalance?apiKey=%v", os.Getenv("MOBIZON_API_KEY"))

	json.Marshal(apiRequest{ApiKey: os.Getenv("MOBIZON_API_KEY")})

	data, _ := json.Marshal(apiRequest{ApiKey: os.Getenv("MOBIZON_API_KEY")})
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))
}

func SendSms(phone string, text string) {

	// эмулируем долгую обработку
	time.Sleep(3 * time.Second)

	if os.Getenv("FAKE_SMS") == "true" {
		log.Printf("Fake sms is enabled. Phone=%v, Text=%v", phone, text)
		return
	}

	var uri *url.URL
	uri, err := url.Parse("https://api.mobizon.ua/service/message/sendSMSMessage")
	if err != nil {
		panic(err)
	}

	parameters := url.Values{}
	parameters.Add("apiKey", os.Getenv("MOBIZON_API_KEY"))
	parameters.Add("recipient", phone)
	parameters.Add("text", text)
	uri.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", uri.String(), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))
}
