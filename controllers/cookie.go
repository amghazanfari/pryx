package controllers

import (
	"net/http"
)

const (
	cookieSession = "session"
)

func setUserSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieSession, // New name, separate from CSRF
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Enable for HTTPS in prod
		SameSite: http.SameSiteLaxMode,
	})
}

func deleteUserSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieSession,
		MaxAge: -1,
	})
}
