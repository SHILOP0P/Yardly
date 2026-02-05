package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ErrorResponse struct{
	Error string `json: "error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}


func QueryInt(r *http.Request, key string, def, min, max int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}