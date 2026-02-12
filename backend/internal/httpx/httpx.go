package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
)

type ErrorResponse struct{
	Error string `json: "error"`
}

// ReadJSON читает JSON из body в dst.
// Делает базовую валидацию: корректный JSON, без лишних данных после объекта.
func ReadJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("empty body")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Проверяем, что в body больше ничего нет (защита от "}{")
	if dec.More() {
		return fmt.Errorf("invalid json")
	}
	// Иногда dec.More() не ловит хвост. Надёжнее так:
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return fmt.Errorf("invalid json")
	}

	return nil
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