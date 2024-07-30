package functions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// User represents the user data from the signup form
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Session store
var store = sessions.NewCookieStore([]byte("secret-key"))

// SignupHandler handles the signup page
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/signup.html")
		return
	}

	if r.Method == http.MethodPost {
		// Get form values
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		passwordConfirmation := r.FormValue("password_confirmation")
		role := r.FormValue("roles")

		// Validate password confirmation
		if password != passwordConfirmation {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Handle user registration logic here (e.g., save to database)
		user := User{
			Name:     name,
			Email:    email,
			Password: password,
			Role:     role,
		}

		// For demonstration, just print the user data
		// In a real application, save the user data to a database
		_ = user

		// Redirect to login
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// LoginHandler handles the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	if r.Method == http.MethodPost {
		// Get form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate credentials (this is a placeholder, replace with actual validation)
		if email == "user@example.com" && password == "password" {
			// Create a new session
			session, err := store.Get(r, "my-session")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			session.Values["authenticated"] = true
			session.Values["email"] = email
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/protected", http.StatusFound)
			return
		}

		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

// LogoutHandler handles the logout page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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

		http.ServeFile(w, r, "templates/logout.html")
	}
}

// ProtectedHandler handles the protected page
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

	w.Write([]byte("Welcome, " + session.Values["email"].(string)))
}

// AuthMiddleware middleware for checking authentication
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
