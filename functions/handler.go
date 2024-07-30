package functions

import (
	"fmt"
	"net/http"
	"text/template"
)

// Index handles the root URL and serves the index.html template
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ServeError(w, "Page not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Learners handles the learners URL and serves the learners.html template
func Learners(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/learners" {
		ServeError(w, "Page not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	tmpl, err := template.ParseFiles("templates/learners.html")
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Tutors handles the tutors URL and serves the tutors.html template
func Tutors(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tutors" {
		ServeError(w, "Page not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	tmpl, err := template.ParseFiles("templates/tutors.html")
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Peers handles the peertopeer URL and serves the index.html template
func Peers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/peertopeer" {
		ServeError(w, "Page not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	tmpl, err := template.ParseFiles("templates/peertopeer.html")
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Admin handles the admin URL and serves the admin.html template
func Admin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		ServeError(w, "Page not Found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		fmt.Println("OK: ", http.StatusOK)
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		ServeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
