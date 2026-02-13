package item

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)

type Handler struct {
	repo Repo
}

func NewHandler(repo Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	f, err := parseListFilter(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	f.Status = []Status{StatusActive, StatusInUse}

	items, err := h.repo.List(r.Context(), *f)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.hydrateImages(r.Context(), items); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	it, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "item not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if it.Status != StatusActive && it.Status != StatusInUse {
		httpx.WriteError(w, http.StatusNotFound, "item not found")
		return
	}

	imgs, err := h.repo.ListImages(r.Context(), it.ID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	it.Images = imgs

	httpx.WriteJSON(w, http.StatusOK, it)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var dto struct {
		Title       string   `json:"title"`
		Mode        DealMode `json:"mode"`
		Description string   `json:"description"`
		Price       int64    `json:"price"`
		Deposit     int64    `json:"deposit"`
		Location    string   `json:"location"`
		Category    string   `json:"category"`
		Images      []struct {
			URL       string `json:"url"`
			SortOrder int    `json:"sort_order"`
		} `json:"images"`
	}
	if err := httpx.ReadJSON(r, &dto); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if dto.Title == "" {
		httpx.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	if !dto.Mode.Valid() {
		httpx.WriteError(w, http.StatusBadRequest, "invalid mode")
		return
	}

	it := Item{
		OwnerID:     ownerID,
		Title:       dto.Title,
		Status:      StatusActive,
		Mode:        dto.Mode,
		Description: dto.Description,
		Price:       dto.Price,
		Deposit:     dto.Deposit,
		Location:    dto.Location,
		Category:    dto.Category,
		Images:      nil,
	}

	if len(dto.Images) > 0 {
		it.Images = make([]ItemImage, 0, len(dto.Images))
		for _, im := range dto.Images {
			if im.URL == "" {
				httpx.WriteError(w, http.StatusBadRequest, "image url is required")
				return
			}
			it.Images = append(it.Images, ItemImage{
				URL:       im.URL,
				SortOrder: im.SortOrder,
			})
		}
	}

	if err := h.repo.Create(r.Context(), &it); err != nil {
		log.Println("item create error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "could not create item")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, it)
}

func (h *Handler) ListMyItems(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	f, err := parseListFilter(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.repo.ListMyItems(r.Context(), ownerID, *f)
	if err != nil {
		log.Println("list my items error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if err := h.hydrateImages(r.Context(), items); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) ListByOwnerPublic(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || ownerID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	f, err := parseListFilter(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.repo.ListByOwnerPublic(r.Context(), ownerID, *f)
	if err != nil {
		log.Println("list owner items error:", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if err := h.hydrateImages(r.Context(), items); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) ListImages(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	imgs, err := h.repo.ListImages(r.Context(), itemID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, imgs)
}

func (h *Handler) AddImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	it, err := h.repo.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "item not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if it.OwnerID != userID {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	url, err := saveUploadedImage(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	im, err := h.repo.AddImage(r.Context(), itemID, url)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, im)
}

func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || itemID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	imageID, err := strconv.ParseInt(r.PathValue("imageId"), 10, 64)
	if err != nil || imageID <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, "invalid image id")
		return
	}

	it, err := h.repo.GetByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "item not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if it.OwnerID != userID {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.repo.DeleteImage(r.Context(), itemID, imageID); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "image not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func parseLimitOffset(r *http.Request) (int, int, error) {
	q := r.URL.Query()

	limit := 20

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return 0, 0, errors.New("invalid limit")
		}
		if n > 100 {
			n = 100
		}
		limit = n
	}
	offset := 0
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return 0, 0, errors.New("invalid offset")
		}
		offset = n
	}

	return limit, offset, nil
}

func parseListFilter(r *http.Request) (*ListFilter, error) {
	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	f := &ListFilter{
		Limit:  limit,
		Offset: offset,
	}

	if v := strings.TrimSpace(q.Get("mode")); v != "" {
		m := DealMode(v)
		if !m.Valid() {
			return nil, errors.New("invalid mode")
		}
		f.Mode = &m
	}

	if v := strings.TrimSpace(q.Get("category")); v != "" {
		f.Category = &v
	}

	if v := strings.TrimSpace(q.Get("location")); v != "" {
		f.Location = &v
	}

	if v := strings.TrimSpace(q.Get("min_price")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			return nil, errors.New("invalid min_price")
		}
		f.MinPrice = &n
	}

	if v := strings.TrimSpace(q.Get("max_price")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			return nil, errors.New("invalid max_price")
		}
		f.MaxPrice = &n
	}

	if f.MinPrice != nil && f.MaxPrice != nil && *f.MinPrice > *f.MaxPrice {
		return nil, errors.New("min_price cannot be greater than max_price")
	}

	return f, nil
}

func saveUploadedImage(r *http.Request) (string, error) {
	const maxFileSize = 10 << 20 // 10MB - битовый сдвиг = 10*2^20

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return "", errors.New("invalid multipart form")
	}

	file, hdr, err := r.FormFile("file")
	if err != nil {
		return "", errors.New("file is required")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(hdr.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", errors.New("only .png, .jpg, .jpeg are allowed")
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", errors.New("failed to read file")
	}
	ct := http.DetectContentType(buf[:n])
	if ct != "image/png" && ct != "image/jpeg" {
		return "", errors.New("only PNG or JPEG content is allowed")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("failed to reset file reader")
	}

	root := os.Getenv("UPLOADS_DIR")
	if strings.TrimSpace(root) == "" {
		root = "uploads"
	}
	itemsDir := filepath.Join(root, "items")
	if err := os.MkdirAll(itemsDir, 0o755); err != nil {
		return "", errors.New("failed to prepare upload directory")
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		return "", errors.New("failed to generate filename")
	}
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(suffix), ext)
	dstPath := filepath.Join(itemsDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New("failed to save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", errors.New("failed to write file")
	}

	return "/uploads/items/" + filename, nil
}

func (h *Handler) hydrateImages(ctx context.Context, items []Item) error {
	for i := range items {
		imgs, err := h.repo.ListImages(ctx, items[i].ID)
		if err != nil {
			return err
		}
		items[i].Images = imgs
	}
	return nil
}
