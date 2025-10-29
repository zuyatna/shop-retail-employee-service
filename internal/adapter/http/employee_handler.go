package http

import (
	"bytes"
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

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

type updatePhotoRequest struct {
	Photo string `json:"photo"` // base64 encoded string
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
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	caller := getCallerRoleFromContext(r)
	items, err := h.empUsecase.FindAll(caller)
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

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	item, err := h.empUsecase.FindByID(callerRole, callerID, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "employee not found"})
		case errors.Is(err, domain.ErrDeleted):
			writeJSON(w, http.StatusGone, map[string]string{"error": "employee has been deleted"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *EmployeeHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")
	id = strings.TrimSuffix(id, "/photo")

	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	employee, err := h.empUsecase.FindByID(callerRole, callerID, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		case errors.Is(err, domain.ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "employee not found"})
		case errors.Is(err, domain.ErrDeleted):
			writeJSON(w, http.StatusGone, map[string]string{"error": "employee has been deleted"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
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

type createEmployeeRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
	Status   string  `json:"status"`
	Address  string  `json:"address"`
	District string  `json:"district"`
	City     string  `json:"city"`
	Province string  `json:"province"`
	Phone    string  `json:"phone"`
	Photo    string  `json:"photo"` // base64 encoded string
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	caller := getCallerRoleFromContext(r)
	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	photo, err := decoderPhotoString(req.Photo)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid photo encoding"})
		return
	}

	employee := &domain.Employee{
		Name:          req.Name,
		Email:         req.Email,
		PasswordHash:  req.Password,
		Role:          domain.Role(req.Role),
		Position:      req.Position,
		Salary:        req.Salary,
		Status:        req.Status,
		Address:       req.Address,
		District:      req.District,
		City:          req.City,
		Province:      req.Province,
		Phone:         req.Phone,
		Photo:         photo,
		PhotoProvided: photo != nil,
	}

	if err := h.empUsecase.Create(caller, employee); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, domain.ErrForbidden):
			status = http.StatusForbidden
		case errors.Is(err, domain.ErrBadRequest):
			status = http.StatusBadRequest
		case errors.Is(err, domain.ErrDuplicate):
			status = http.StatusConflict
		case errors.Is(err, domain.ErrPhotoTooLarge):
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Println("Employee created with ID:", employee.ID)
	writeJSON(w, http.StatusCreated, map[string]string{"id": employee.ID})
}

type updateEmployeeRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password *string `json:"password"`
	Role     string  `json:"role"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
	Status   string  `json:"status"`
	Address  string  `json:"address"`
	District string  `json:"district"`
	City     string  `json:"city"`
	Province string  `json:"province"`
	Phone    string  `json:"phone"`
	Photo    *string `json:"photo"` // base64 encoded string
}

func (h *EmployeeHandler) PutPhoto(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/employee/"), "/photo")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	var req updatePhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	data, mime, err := decodeAndDetectMIME(req.Photo)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid photo encoding"})
		return
	}

	if err := h.empUsecase.UpdatePhoto(callerRole, callerID, id, data, mime); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrDeleted) {
			status = http.StatusGone
		} else if errors.Is(err, domain.ErrPhotoTooLarge) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Printf("Employee with ID %s photo updated\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee photo updated"})
}

func (h *EmployeeHandler) PutPhotoMultipart(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/employee/"), "/photo")
	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	if err := r.ParseMultipartForm(6 << 20); err != nil { // 6 MB limit form
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form data"})
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to read photo file"})
		return
	}
	defer file.Close()

	// read max 5 MB + 1 byte to check size
	const maxPhotoSize = 5 * 1024 * 1024
	buf := bytes.NewBuffer(nil)

	if _, err := io.CopyN(buf, file, maxPhotoSize+1); err != nil && err != io.EOF {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to read photo file"})
		return
	}

	if buf.Len() > maxPhotoSize {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "photo size exceeds the limit"})
		return
	}

	data := buf.Bytes()
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

	if err := h.empUsecase.UpdatePhoto(callerRole, callerID, id, data, mime); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrDeleted) {
			status = http.StatusGone
		} else if errors.Is(err, domain.ErrPhotoTooLarge) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Printf("Employee with ID %s photo updated\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee photo updated"})

}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")

	callerRole := getCallerRoleFromContext(r)
	callerID := getCallerIDFromContext(r)

	var req updateEmployeeRequest
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

	if err := h.empUsecase.Update(callerRole, callerID, employee); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrDeleted) {
			status = http.StatusGone
		} else if errors.Is(err, domain.ErrPhotoTooLarge) {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Printf("Employee with ID %s updated\n", id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "employee updated"})
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")
	caller := getCallerRoleFromContext(r)
	if err := h.empUsecase.Delete(caller, id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrDeleted) {
			status = http.StatusGone
		} else if errors.Is(err, domain.ErrForbidden) {
			status = http.StatusForbidden
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Printf("Employee with ID %s deleted\n", id)
	w.WriteHeader(http.StatusNoContent)
}

func decodeAndDetectMIME(s string) ([]byte, string, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil, "", errors.New("empty photo")
	}

	var hintedMIME string
	if idx := strings.Index(trimmed, ","); idx != -1 && strings.Contains(strings.ToLower(trimmed[:idx]), "base64") {
		header := trimmed[:idx]
		trimmed = trimmed[idx+1:]
		if p := strings.Index(header, ":"); p != -1 {
			if q := strings.Index(header[p+1:], ";"); q != -1 {
				hintedMIME = header[p+1 : p+1+q]
			}
		}
	}

	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil || len(data) == 0 {
		return nil, "", errors.New("invalid base64 encoding")
	}

	sniff := http.DetectContentType(data[:min(512, len(data))])
	mime := strings.ToLower(strings.TrimSpace(hintedMIME))
	if mime == "" {
		mime = strings.ToLower(sniff)
	}

	switch mime {
	case "image/jpeg", "image/jpg":
		mime = "image/jpeg"
		if !(len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF) {
			return nil, "", errors.New("invalid JPEG/JPG image data signature")
		}
	case "image/png":
		sig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if !(len(data) >= 8 && string(data[:8]) == string(sig)) {
			return nil, "", errors.New("invalid PNG image data signature")
		}
	default:
		return nil, "", errors.New("unsupported image MIME type")
	}
	return data, mime, nil
}
