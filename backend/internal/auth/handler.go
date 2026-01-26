package auth

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct{
	jwt *JWT
}

func NewHandler(jwt *JWT) *Handler{
	return &Handler{jwt: jwt}
}

type loginRequest struct{
	UserID int64 `json:"user_id"`
}

type loginResponse struct{
	AccessToken string `json:"access_token"`
	TokenType 	 string `json:"token_type"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request){
	var req loginRequest
	if err :=json.NewDecoder(r.Body).Decode(&req); err!=nil||req.UserID<=0{
		log.Println("Invalid user_id or failed to decode request body:", err)
		httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	tok, err := h.jwt.Mint(req.UserID)
	if err != nil {
		log.Println("Error minting token:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint token")
		return
	}
	log.Println("Generated token for user:", req.UserID)
	httpx.WriteJSON(w, http.StatusOK, loginResponse{
		AccessToken: tok,
		TokenType:   "Bearer",
})

}
