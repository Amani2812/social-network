package auth

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const SessionCookieName = "session_id"

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

var sessions = make(map[string]*Session)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CreateSession(userID int) (*Session, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sessions[sessionID] = session
	return session, nil
}

func GetSession(sessionID string) (*Session, bool) {
	session, exists := sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		if exists {
			delete(sessions, sessionID)
		}
		return nil, false
	}
	return session, true
}

func DeleteSession(sessionID string) {
	delete(sessions, sessionID)
}

func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func GetUserIDFromRequest(r *http.Request) (int, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0, false
	}

	session, exists := GetSession(cookie.Value)
	if !exists {
		return 0, false
	}

	return session.UserID, true
}
