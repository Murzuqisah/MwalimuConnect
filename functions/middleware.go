package functions

import (
	"context"
	"net/http"
	"github.com/google/uuid"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := InitSession(w, r)
		if session == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		next.ServeHTTP(w, r)
	})
}

func InitSession(w http.ResponseWriter, r *http.Request) *Session {
	cookie, err := r.Cookie("session")
	if err != nil {
		// Create a new session and set the cookie
		session := &Session{ID: GenerateSessionID()}
		cookie = &http.Cookie{
			Name:  "session",
			Value: session.ID,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		return session
	}

	// Retrieve the session from the store using the session ID
	session, err := RetrieveSession(cookie.Value)
	if err != nil {
		return nil
	}
	return session
}

func GenerateSessionID() string {
	id := uuid.New()
	return id.String()
}

func RetrieveSession(sessionID string) (*Session, error) {
	// TO DO: implement session retrieval logic here
	// For example, you can use a database or a cache to store sessions
	// Return the session object or an error if not found
	return &Session{ID: sessionID}, nil
}

type Session struct {
	ID string
	// Add other session fields as needed
}
