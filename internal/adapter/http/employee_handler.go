package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/zuyatna/shop-retail-employee-service/internal/model"
	domain "github.com/zuyatna/shop-retail-employee-service/internal/model"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

// read max 5 MB + 1 byte to check size
const maxPhotoSize = 5 * 1024 * 1024

var bufPool = sync.Pool{
	New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) },
}

type EmployeeHandler struct {
	empUsecase *usecase.EmployeeUsecase
}

func NewEmployeeHandler(empUsecase *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		empUsecase: empUsecase,
	}
}

func decoderPhotoString(s string) ([]byte, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil, nil
	}
	if idx := strings.Index(trimmed, ","); idx != -1 {
		prefix := trimmed[:idx]
		if strings.Contains(strings.ToLower(prefix), "base64") {
			trimmed = trimmed[idx+1:]
		}
	}

	rdr := base64.NewDecoder(base64.StdEncoding, strings.NewReader(trimmed))
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufPool.Put(b)

	if _, err := io.Copy(b, rdr); err != nil {
		return nil, err
	}

	out := append([]byte(nil), b.Bytes()...) // copy before returning if pool reuse concern
	return out, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErrorJSON(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrDeleted):
		status = http.StatusGone
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (h *EmployeeHandler) List(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	caller := getCallerRoleFromContext(r)

	items, err := h.empUsecase.FindAll(ctx, caller)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrForbidden) {
			status = http.StatusForbidden
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *EmployeeHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	item, err := h.empUsecase.FindByID(ctx, callerRole, callerID, id)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *EmployeeHandler) GetMe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	if strings.TrimSpace(callerID) == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	employee, err := h.empUsecase.FindByID(ctx, callerRole, callerID, callerID)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	employeeNew := model.EmployeeResponse{
		ID:    employee.ID,
		Name:  employee.Name,
		Email: employee.Email,
		Role:  employee.Role,
	}

	writeJSON(w, http.StatusOK, &employeeNew)
}

func (h *EmployeeHandler) GetPhoto(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/photo/")

	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	employee, err := h.empUsecase.FindByID(ctx, callerRole, callerID, id)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	if len(employee.Photo) == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "photo not found"})
		return
	}

	sum := sha256.Sum256(employee.Photo)
	etag := `W/"` + hex.EncodeToString(sum[:]) + `"`
	if match := r.Header.Get("If-None-Match"); match != "" && match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	mime := employee.PhotoMIME
	if mime == "" {
		mime = "application/octet-stream" // default MIME type
	}

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(employee.Photo)))
	w.Header().Set("Cache-Control", "private, max-age=86400") // cache for 1 day
	w.Header().Set("ETag", etag)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(employee.Photo)
}

func (h *EmployeeHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	caller := getCallerRoleFromContext(r)
	var req model.EmployeeRequest
	var pwHash string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password is required"})
		return
	} else {
		pwHash = strings.TrimSpace(*req.Password)
	}

	employee := &domain.Employee{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: pwHash,
		Role:         domain.Role(req.Role),
		Position:     req.Position,
		Salary:       req.Salary,
		Status:       req.Status,
		Address:      req.Address,
		District:     req.District,
		City:         req.City,
		Province:     req.Province,
		Phone:        req.Phone,
	}

	if err := h.empUsecase.Create(ctx, caller, employee); err != nil {
		writeErrorJSON(w, err)
		return
	}
	log.Println("Employee created with ID:", employee.ID)
	writeJSON(w, http.StatusCreated, map[string]string{"id": employee.ID})
}

func (h *EmployeeHandler) PutPhotoMultipart(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/photo/")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	r.Body = http.MaxBytesReader(w, r.Body, 6<<20) // 6 MB max
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to parse multipart form"})
		return
	}

	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bufPool.Put(b)

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read photo file"})
		return
	}
	defer file.Close()

	if _, err := io.CopyN(b, file, maxPhotoSize+1); err != nil && err != io.EOF {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read photo file"})
		return
	}

	if b.Len() > maxPhotoSize {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "photo size exceeds the limit"})
		return
	}

	// copy data from buffer pool to new slice
	data := make([]byte, b.Len())
	copy(data, b.Bytes())

	// detect MIME type
	sniff := http.DetectContentType(data[:min(512, len(data))])
	mime := strings.ToLower(sniff)

	switch mime {
	case "image/jpeg", "image/jpg":
		mime = "image/jpeg"
	case "image/png":
		mime = "image/png"
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported image MIME type"})
		return
	}

	// validate signature
	switch mime {
	case "image/jpeg", "image/jpg":
		mime = "image/jpeg"
		if !(len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JPEG/JPG image data signature"})
			return
		}
	case "image/png":
		sig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if !(len(data) >= 8 && string(data[:8]) == string(sig)) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid PNG image data signature"})
			return
		}
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported image MIME type"})
		return
	}

	_ = header // to avoid unused variable warning

	if err := h.empUsecase.UpdatePhoto(ctx, callerRole, callerID, id, data, mime); err != nil {
		writeErrorJSON(w, err)
		return
	}
	log.Printf("Employee with ID %s photo updated\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee photo updated"})
}

func (h EmployeeHandler) DeletePhoto(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/photo/")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	if err := h.empUsecase.UpdatePhoto(ctx, callerRole, callerID, id, nil, ""); err != nil {
		writeErrorJSON(w, err)
		return
	}
	log.Printf("Employee with ID %s photo deleted\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee photo deleted"})
}

func (h *EmployeeHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")

	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	var req model.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	employee := &domain.Employee{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Role:     domain.Role(req.Role),
		Position: req.Position,
		Salary:   req.Salary,
		Status:   req.Status,
		Address:  req.Address,
		District: req.District,
		City:     req.City,
		Province: req.Province,
		Phone:    req.Phone,
	}
	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password != "" {
			employee.PasswordHash = password
		}
	}

	if req.Photo != nil {
		photo, err := decoderPhotoString(*req.Photo)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid photo encoding"})
			return
		}
		if len(photo) == 0 {
			employee.Photo = nil
		} else {
			employee.Photo = photo
		}
		employee.PhotoProvided = true
	}

	if err := h.empUsecase.Update(ctx, callerRole, callerID, employee); err != nil {
		writeErrorJSON(w, err)
		return
	}
	log.Printf("Employee with ID %s updated\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee updated"})
}

func (h *EmployeeHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")
	caller := getCallerRoleFromContext(r)
	if err := h.empUsecase.Delete(ctx, caller, id); err != nil {
		writeErrorJSON(w, err)
		return
	}
	log.Printf("Employee with ID %s deleted\n", id)
	w.WriteHeader(http.StatusNoContent)
}
