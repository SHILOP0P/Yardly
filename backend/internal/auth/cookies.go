package auth

import (
	"net/http"
	"time"
)

const refreshCookieName = "refresh_token"

// ВАЖНО:
// - Secure=true только на https
// - SameSite=Lax обычно ок для SPA (и защищает от части CSRF)
// - Path лучше ограничить, но чтобы logout тоже видел cookie, ставим "/api/auth"
func setRefreshCookie(w http.ResponseWriter, token string, expiresAt time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/api/auth",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearRefreshCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/api/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func readRefreshCookie(r *http.Request) (string, error) {
	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		return "", err
	}
	if c.Value == "" {
		return "", http.ErrNoCookie
	}
	return c.Value, nil
}
