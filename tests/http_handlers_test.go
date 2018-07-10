package tests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/JILeXanDR/golang/http_handlers"
	"strings"
	"github.com/JILeXanDR/golang/app"
	"encoding/json"
)

func assertEqual(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Fatalf("got %v != expected %v", got, expected)
	}
}

func createResponse(req *http.Request, fn func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder) {

	app.CreateTest()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

func TestHomePageOpened(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := createResponse(req, http_handlers.HomePageHandler)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `<div id="app"></div>`

	if content := res.Body.String(); !strings.Contains(content, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
	}
}

func TestGetBalance(t *testing.T) {

	req, _ := http.NewRequest("GET", "/balance", strings.NewReader(`{"user": 10}`))
	res := createResponse(req, http_handlers.GetBalanceHandler)

	decoder := json.NewDecoder(res.Body)
	var data http_handlers.BalanceResponse
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	assertEqual(t, int(data.Balance), 1000)
}
