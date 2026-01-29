package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"log"
	"time"

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
	refreshTok, err := mintRefreshTocken()
	if err!=nil{
		log.Println("Error minting refresh token:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint refresh token")
		return
	}

	expiresAt := time.Now().Add(h.refreshTTL)
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

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request){
	var req refreshRequest
	if err:= json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken ==""{
		httpx.WriteError(w, http.StatusBadRequest, "invalid refresh_token")
		return
	}
	
	now:=time.Now().UTC()

	userID, err := h.refresh.Consume(r.Context(), req.RefreshToken, now)
	if err!=nil{
		httpx.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessTok, err := h.jwt.Mint(userID)
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint token")
		return
	}

	newRefresh, err := mintRefreshTocken()
	if err!=nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint refresh token")
		return
	}

	expiresAt := now.Add(h.refreshTTL)
	if err := h.refresh.Create(r.Context(), userID, newRefresh, expiresAt); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "could not save refresh token")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, refreshResponse{
		AccessToken:  accessTok,
		RefreshToken: newRefresh,
		TokenType:    "Bearer",
	})
}












func mintRefreshTocken()(string, error){
	b:=make([]byte, 32)
	if _, err :=rand.Read(b); err != nil{
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}