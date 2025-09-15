package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Context key
type key int

const userIDKey key = 0

// Session struct
type Session struct {
	Username string
	Expiry   time.Time
}

// Sessions stores sessionID -> Session
var Sessions sync.Map

// Session duration
const sessionDuration = 15 * time.Minute

// CreateSession generates a new sessionID for a user and sets it in the cookie
func CreateSession(w http.ResponseWriter, username string) string {
	b := make([]byte, 16)
	rand.Read(b)
	sessionID := fmt.Sprintf("%x", b)

	Sessions.Store(sessionID, Session{
		Username: username,
		Expiry:   time.Now().Add(sessionDuration),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		Expires:  time.Now().Add(sessionDuration),
	})

	return sessionID
}

// CookieAuthMiddleware checks for a valid session cookie
func CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_id")
		if err != nil || c.Value == "" {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		val, ok := Sessions.Load(c.Value)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		session := val.(Session)
		if session.Expiry.Before(time.Now()) {
			// session expired
			Sessions.Delete(c.Value)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Optional: refresh session expiry on each request
		session.Expiry = time.Now().Add(sessionDuration)
		Sessions.Store(c.Value, session)

		ctx := context.WithValue(r.Context(), userIDKey, session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext retrieves the username from context
func UserIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(userIDKey).(string)
	return uid, ok
}

// StartSessionCleanup cleanup old sessions periodically
func StartSessionCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			now := time.Now()
			Sessions.Range(func(key, value any) bool {
				sess := value.(Session)
				if sess.Expiry.Before(now) {
					Sessions.Delete(key)
				}
				return true
			})
		}
	}()
}
