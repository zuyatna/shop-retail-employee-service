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
	empUsecase *usecase.EmployeeUsecase
}

func NewEmployeeHandler(empUsecase *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		empUsecase: empUsecase,
	}
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
	callerID := getCalledIDFromContext(r)

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
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	caller := getCallerRoleFromContext(r)
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
		Address:      req.Address,
		District:     req.District,
		City:         req.City,
		Province:     req.Province,
		Phone:        req.Phone,
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
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/employee/")

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

	if err := h.empUsecase.Update(employee); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrBadRequest) {
			status = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, domain.ErrDeleted) {
			status = http.StatusGone
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
