package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"log"
	"time"
	stdErrors "errors"

	"crypto/rand"
	"encoding/base64"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct {
	jwt        *JWT
	refresh    *RefreshRepo
	refreshTTL time.Duration
	users UserAuthenticator
}

func NewHandler(jwt *JWT, refresh *RefreshRepo, refreshTTL time.Duration, users UserAuthenticator) *Handler {
	return &Handler{jwt: jwt, refresh: refresh, refreshTTL: refreshTTL, users: users}
}


type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type UserAuthenticator interface {
	Authenticate(ctx context.Context, email, password string) (int64, error)
}



type refreshRequest struct{
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}


func (h *Handler) Login(w http.ResponseWriter, r *http.Request){
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	if req.Email == "" || req.Password == "" {
		httpx.WriteError(w, http.StatusBadRequest, "email and password required")
		return
	}

	userID, err := h.users.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		// 401 чтобы не палить, существует ли email
		httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	tok, err := h.jwt.Mint(userID)
	if err != nil {
		log.Println("Error minting token:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint token")
		return
	}
	refreshTok, err := mintRefreshToken()
	if err!=nil{
		log.Println("Error minting refresh token:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint refresh token")
		return
	}

	expiresAt := time.Now().UTC().Add(h.refreshTTL)
	if err :=h.refresh.Create(r.Context(), userID, refreshTok, expiresAt); err!= nil{
		log.Println("Error saving refresh token:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not save refresh token")
		return
	}

	log.Println("Generated token for user:", userID)
	httpx.WriteJSON(w, http.StatusOK, loginResponse{
		AccessToken: tok,
		RefreshToken: refreshTok,
		TokenType:   "Bearer",
	})
}


const refreshGrace = 5 * time.Second


func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request){
	var req refreshRequest
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken ==""{
		httpx.WriteError(w, http.StatusBadRequest, "invalid refresh_token")
		return
	}
	
	now:=time.Now().UTC()
	
	expiresAt := now.Add(h.refreshTTL)

	newRefresh, err := mintRefreshToken()
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint refresh token")
		return
	}

	userID, err := h.refresh.Rotate(r.Context(), req.RefreshToken, newRefresh, expiresAt, now, refreshGrace)
	if err != nil {
		switch {
		case stdErrors.Is(err, ErrInvalidRefresh):
			httpx.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		case stdErrors.Is(err, ErrRefreshAlreadyRotated):
			httpx.WriteError(w, http.StatusUnauthorized, "stale refresh token")
		case stdErrors.Is(err, ErrRefreshReuse):
			httpx.WriteError(w, http.StatusUnauthorized, "refresh reuse detected")
		default:
			log.Println(err)
			httpx.WriteError(w, http.StatusInternalServerError, "could not rotate refresh token")
		}
		return
	}

	accessTok, err := h.jwt.Mint(userID)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint token")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, refreshResponse{
		AccessToken:  accessTok,
		RefreshToken: newRefresh,
		TokenType:    "Bearer",
	})
}


func (h *Handler) Logout(w http.ResponseWriter, r *http.Request){
	var req refreshRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil || req.RefreshToken == ""{
		httpx.WriteError(w, http.StatusBadRequest, "invalid refresh_token")
		return
	}
	now := time.Now().UTC()

	if err := h.refresh.Revoke(r.Context(), req.RefreshToken, now); err!= nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not revoke refresh token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request){
	userID, ok := UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
	}

	now := time.Now().UTC()

	if err:=h.refresh.RevokeAllForUser(r.Context(),userID, now); err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not revoke sessions")
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}











func mintRefreshToken()(string, error){
	b:=make([]byte, 32)
	if _, err :=rand.Read(b); err != nil{
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}