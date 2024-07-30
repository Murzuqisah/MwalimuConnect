package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"mwalimu/functions"
)

func main() {
	http.HandleFunc("/", servePage)
	// Initialize blockchain with genesis block
	blockchain := functions.Blockchain{}
	genesisBlock := functions.Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte("Genesis Block"),
		PreviousHash: []byte{},
	}
	genesisBlock.Hash = functions.CalculateHash(genesisBlock)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	// Serve static files (CSS, JS)
	http.HandleFunc("/scripts/", fileServer)
	http.HandleFunc("/static/", fileServer)

	// Serve dynamic HTML templates
	//http.HandleFunc("/", servePage)

	// // Handle routes
	http.HandleFunc("/forgot-password", functions.ServeForgotPasswordForm)
	http.HandleFunc("/handle-forgot-password", functions.HandleForgotPassword)
	http.HandleFunc("/reset-password", functions.ServeResetPasswordForm)
	http.HandleFunc("/handle-reset-password", functions.HandleResetPassword)

	// serve signup
	http.HandleFunc("/signup", functions.SignupHandler)

	// Start the server
	log.Printf("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Serve static files (CSS, JS)
func fileServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.URL.Path[len("/scripts/"):])
}

// Serve HTML pages based on referral
func servePage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if page == "" || page == "/" {
		page = "index.html"
	}

	// Define template path
	tmplPath := "./templates/" + page

	// Parse and execute the template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

// The handler function for signup route
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		functions.SignupHandler(w, r)
	default:
		code := http.StatusNotFound
		codefin := strconv.Itoa(code)
		fmt.Printf(codefin, "Page not found")
		os.Exit(1)
	}
}
