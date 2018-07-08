package http_handlers

import (
	"net/http"
	"io/ioutil"
)

// отображение главной страницы
func HomePageHandler(w http.ResponseWriter, r *http.Request) {

	html, err := ioutil.ReadFile("./public/index.html")
	if err != nil {
		panic(err)
	}

	_, err = w.Write(html)
	if err != nil {
		panic(err)
	}
}
