package functions

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"gopkg.in/mail.v2"
)

var users = make(map[string]string) // email -> token (in a real app, store in database)

func ServeForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/forgot_password.html"))
	tmpl.Execute(w, nil)
}

func HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Check if the email exists in your user database
	// For simplicity, assume user exists
	// Generate a unique token
	token := uuid.New().String()

	// Store the token with the email in your database or map
	users[email] = token

	// Send password reset email with token link
	err := sendPasswordResetEmail(email, token)
	if err != nil {
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}

	// Display confirmation message to the user
	fmt.Fprintf(w, "Password reset instructions sent to %s", email)
}

func sendPasswordResetEmail(email, token string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "your-email@example.com") // Replace with your email
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset Instructions")

	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	m.SetBody("text/html", fmt.Sprintf("Click <a href=\"%s\">here</a> to reset your password.", resetURL))

	d := mail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password") // Replace with your SMTP settings
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}
	fmt.Println("Password reset email sent to:", email)
	return nil
}
