package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct {
	repo Repo
	jwt *auth.JWT
}

func NewHandler(repo Repo, jwt *auth.JWT) *Handler{
	return &Handler{repo: repo, jwt: jwt}
}

type registerRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
}

type registerResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	FirstName string `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request){
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	req.FirstName = strings.TrimSpace(req.FirstName)
	if req.LastName != nil{
		trimmed := strings.TrimSpace(*req.LastName)
		if trimmed == ""{
			req.LastName = nil
		} else{
			req.LastName = &trimmed
		}
	}
	if req.Email == "" || !strings.Contains(req.Email, "@"){
		httpx.WriteError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(req.Password)<8{
		httpx.WriteError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if req.FirstName == "" {
		httpx.WriteError(w, http.StatusBadRequest, "first_name is required")
		return
	}

	if _, err := h.repo.GetByEmail(r.Context(), req.Email); err == nil{
		httpx.WriteError(w, http.StatusConflict, "email already taken")
		return
	} else if err != ErrNotFound{
		httpx.WriteError(w, http.StatusInternalServerError, "could not check email")
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil{
		httpx.WriteError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	u:= User{
		Email:        req.Email,
		PasswordHash: hash,
		Role:         RoleUser,
	}
	p := Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.repo.CreateWithProfile(r.Context(), &u, &p); err != nil {
		if errors.Is(err, ErrEmailTaken) {
			httpx.WriteError(w, http.StatusConflict, "email already taken")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, registerResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request){
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	pass := strings.TrimSpace(req.Password)

	if email == "" || !strings.Contains(email, "@") {
		httpx.WriteError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if pass == "" {
		httpx.WriteError(w, http.StatusBadRequest, "password is required")
		return
	}

	u, err := h.repo.GetByEmail(r.Context(), email)
	if err != nil{
		httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !CheckPassword(u.PasswordHash, pass){
		httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tok, err := h.jwt.Mint(u.ID, "", true)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "could not mint token")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, loginResponse{
		AccessToken: tok,
		TokenType:   "Bearer",
	})
}

type meResponse struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	Role      Role    `json:"role"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request){
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok{
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, p, err := h.repo.GetByID(r.Context(), userID)
	if err!=nil{
		httpx.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, meResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	})
}
