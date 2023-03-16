package pkg

import (
	"log"
	"net/http"
	"text/template"
)

type errorData struct {
	StatusText string
	StatusCode int
}

func ErrorHandler(w http.ResponseWriter, status int) {
	data := errorData{
		StatusText: http.StatusText(status),
		StatusCode: status,
	}
	tmpl, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		log.Println("asdasdsa")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
