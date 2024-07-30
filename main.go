package main

import (
	"log"
	"net/http"

	"mwalimu/functions"
)

func main() {
	// Serve static files from the "templates" folder
	http.Handle("/", http.FileServer(http.Dir("./templates")))

	// Register handler functions
	http.HandleFunc("/signup", functions.SignupHandler)
	http.HandleFunc("/login", functions.LoginHandler)
	http.HandleFunc("/logout", functions.LogoutHandler)
	http.Handle("/protected", functions.AuthMiddleware(http.HandlerFunc(functions.ProtectedHandler)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
