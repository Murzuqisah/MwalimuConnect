package functions

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// Session store
var store = sessions.NewCookieStore([]byte("secret-key"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new session
	session, err := store.Get(r, "my-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values
	session.Values["username"] = "john"
	session.Values["authenticated"] = true

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/protected", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current session
	session, err := store.Get(r, "my-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Destroy the session
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current session
	session, err := store.Get(r, "my-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is authenticated
	if !session.Values["authenticated"].(bool) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, "Welcome, "+session.Values["username"].(string))
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "my-session")
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return

		}

		if !session.Values["authenticated"].(bool) {

			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return

		}

		next.ServeHTTP(w, r)
	})
}
