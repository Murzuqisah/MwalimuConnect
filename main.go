package main

import (
	"net/http"
	"time"

	"mwalimu/functions"
)

func main() {
	http.HandleFunc("/login", functions.LoginHandler)
	http.HandleFunc("/logout", functions.LogoutHandler)
	http.Handle("/protected", functions.AuthMiddleware(http.HandlerFunc(functions.ProtectedHandler)))

	blockchain := functions.Blockchain{}
	genesisBlock := functions.Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte("Genesis Block"),
		PreviousHash: []byte{},
	}
	genesisBlock.Hash = functions.CalculateHash(genesisBlock)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	block1Data := []byte("")

	// Serve static files (CSS, JS)
	http.HandleFunc("/scripts/", fileServer)
	http.HandleFunc("/static/", fileServer)

	// Serve dynamic HTML templates
	http.HandleFunc("/", servePage)

	// Start the server
	http.ListenAndServe(":8080", nil)
}

// Serve static files (CSS, JS)
func fileServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.URL.Path[1:])
}

// Serve HTML pages based on referral
func servePage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if page == "" || page == "/" {
		page = "admin.html" // Default page
	}

	// Define template path
	tmplPath := "./templates/" + page + ".html"

	// Parse and execute the template
	tmpl, err := functions.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, nil)
}
