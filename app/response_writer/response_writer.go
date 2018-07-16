package response_writer

import (
	"net/http"
	"encoding/json"
	"os"
	"errors"
	"log"
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
	log.Println(err)
	if os.Getenv("ENV") == "dev" {
		HandleError(w, err)
	} else {
		HandleError(w, errors.New("Internal Server Error"))
	}
}

func JsonMessageResponse(w http.ResponseWriter, message string, code int) {
	JsonResponse(w, map[string]string{"message": message}, code)
}

func ValidationError(w http.ResponseWriter, message string) {
	JsonMessageResponse(w, message, 422)
}
