package http_handlers

import (
	"net/http"
	"io/ioutil"
	"os"
)

// отображение главной страницы
func HomePageHandler(w http.ResponseWriter, r *http.Request) {

	file := os.Getenv("ROOT_DIR") + "/public/index.html"

	html, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	_, err = w.Write(html)
	if err != nil {
		panic(err)
	}
}
