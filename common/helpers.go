package common

import (
	"net/http"
	"encoding/json"
	"errors"
	"os"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func HandleError(w http.ResponseWriter, err error) {
	JsonMessageResponse(w, err.Error(), 500)
}

func InternalServerError(w http.ResponseWriter, err error) {
	if os.Getenv("ENV") == "dev" {
		panic(err)
		HandleError(w, err)
	} else {
		HandleError(w, errors.New("Internal Server Error"))
	}
}

func ValidationError(w http.ResponseWriter, message string) {
	JsonMessageResponse(w, message, 422)
}

func JsonMessageResponse(w http.ResponseWriter, message string, code int) {
	JsonResponse(w, map[string]string{"message": message}, code)
}
