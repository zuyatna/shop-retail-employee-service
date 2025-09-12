package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
)

type EmployeeHandler struct {
	svc *usecase.EmployeeUsecase
}

func NewEmployeeHandler(svc *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		svc: svc,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employees/")
	item, err := h.svc.FindByID(id)
	if err != nil {
		switch {
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

type createEmployeeRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
	Status   string  `json:"status"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	employee := &domain.Employee{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         domain.Role(req.Role),
		Position:     req.Position,
		Salary:       req.Salary,
		Status:       req.Status,
	}

	if err := h.svc.Create(employee); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrDuplicate) {
			status = http.StatusConflict
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	log.Println("Employee created with ID:", employee.ID)
	writeJSON(w, http.StatusCreated, map[string]string{"id": employee.ID})
}
