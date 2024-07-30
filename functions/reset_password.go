package functions

import (
	"fmt"
	"net/http"
	"text/template"
)

func ServeResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/reset.html"))
	tmpl.Execute(w, token)
}

func HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Validate the token (check against stored tokens)
	// For simplicity, assume token is valid

	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if newPassword == "" || confirmPassword == "" {
		http.Error(w, "New password and confirm password are required", http.StatusBadRequest)
		return
	}

	if newPassword != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Reset password logic here (update in database, etc.)
	// For demonstration, we'll just print the new password
	fmt.Println("New password:", newPassword)

	// Redirect to a success page or display a success message
	http.Redirect(w, r, "/password-reset-success", http.StatusSeeOther)
}
