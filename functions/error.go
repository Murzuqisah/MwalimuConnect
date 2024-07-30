package functions

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Message struct {
	Code       string
	ErrMessage string
}

func ServeError(w http.ResponseWriter, errVal string, statusCode int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("error parsing files: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	code := strconv.Itoa(statusCode)

	data := Message{
		Code:       code,
		ErrMessage: errVal,
	}

	w.WriteHeader(statusCode) // Set the HTTP status code

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
