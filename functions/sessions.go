package functions

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

// Session store
var store = sessions.NewCookieStore([]byte("secret-key"))

func SetSessionValues(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, err := store.Get(r, "my-session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(r, w)
}

func GetSessionValue(r *http.Request, key string) interface{} {
	session, _ := store.Get(r, "my-session")
	return session.Values[key]
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "my-session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
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

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the signup form
		tmpl := `
		<!DOCTYPE html>
		<html>
		<head><title>Sign Up</title></head>
		<body>
			<h1>Sign Up</h1>
			<form method="post" action="/signup">
				<label>Name:</label>
				<input type="text" name="name" required><br>
				<label>Email:</label>
				<input type="text" name="email" required><br>
				<label>Password:</label>
				<input type="password" name="password" required><br>
				<label for="password_confirmation">Repeat Password</label>
            	<input type="password" id="password_confirmation" name="password_confirmation" required>
				<button type="submit">Sign Up</button>
			</form>
		</body>
		</html>
		`
		t := template.New("signup")
		t, _ = t.Parse(tmpl)
		t.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Handle form submission
		username := r.FormValue("username")
		password := r.FormValue("password")

		if _, exists := users[username]; exists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		// Store the user credentials (in a real app, hash the password)
		users[username] = password

		// Redirect to login page
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
