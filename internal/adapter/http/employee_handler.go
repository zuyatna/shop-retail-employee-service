package adapterhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zuyatna/shop-retail-employee-service/internal/domain"
	"github.com/zuyatna/shop-retail-employee-service/internal/dto/employee"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
)

type EmployeeHandler struct {
	usecase *usecase.EmployeeUsecase
}

func NewEmployeeHandler(uc *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{
		usecase: uc,
	}
}

func (h *EmployeeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(UserClaimsKey).(*jwtutil.Claims)
	if !ok || claims == nil {
		WriteErrorJSON(w, http.StatusUnauthorized, nil, "failed to get user context")
		return
	}

	ctx := r.Context()

	getByID, err := h.usecase.GetByID(ctx, claims.UserID)
	if err != nil {
		if err.Error() == usecase.EmployeeNotFoundError {
			WriteErrorJSON(w, http.StatusNotFound, err, "employee not found")
			return
		}

		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to retrieve employee")
		return
	}

	// Convert domain entity to response DTO
	resp := usecase.FromDomain(getByID)
	WriteJSON(w, http.StatusOK, resp, "employee retrieved successfully")
}

func (h *EmployeeHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req employee.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	ctx := r.Context()

	id, err := h.usecase.Register(ctx, req)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"id": id}, "employee registered successfully")
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Assume we get the ID from the URL path, e.g., /employees/{id}
	id := r.URL.Path[len("/employees/"):]

	ctx := r.Context()

	getByID, err := h.usecase.GetByID(ctx, id)
	if err != nil {
		if err.Error() == usecase.EmployeeNotFoundError {
			WriteErrorJSON(w, http.StatusNotFound, err, "getByID not found")
			return
		}

		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to retrieve getByID")
		return
	}

	// Convert domain entity to response DTO
	resp := usecase.FromDomain(getByID)
	WriteJSON(w, http.StatusOK, resp, "getByID retrieved successfully")
}

func (h *EmployeeHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	ctx := r.Context()

	getByEmail, err := h.usecase.GetByEmail(ctx, email)
	if err != nil {
		if err.Error() == usecase.EmployeeNotFoundError {
			WriteErrorJSON(w, http.StatusNotFound, err, "getByEmail not found")
			return
		}
		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to retrieve getByEmail")
		return
	}

	// Convert domain entity to response DTO
	resp := usecase.FromDomain(getByEmail)
	WriteJSON(w, http.StatusOK, resp, "getByEmail retrieved successfully")
}

func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employees, err := h.usecase.GetAll(ctx)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to retrieve employees")
		return
	}

	// Convert domain entities to response DTOs
	var resp []*usecase.EmployeeResponse
	for _, emp := range employees {
		resp = append(resp, usecase.FromDomain(emp))
	}

	// Ensure resp is an empty array if no employees found
	if resp == nil {
		resp = []*usecase.EmployeeResponse{}
	}

	WriteJSON(w, http.StatusOK, resp, "employees retrieved successfully")
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Assume we get the ID from the URL path, e.g., /employees/{id}
	id := r.URL.Path[len("/employees/"):]
	if id == "" {
		WriteErrorJSON(w, http.StatusBadRequest, nil, "employee ID is required")
		return
	}

	claims, ok := r.Context().Value(UserClaimsKey).(*jwtutil.Claims)
	if ok && claims.Role == string(domain.RoleStaff) {
		if claims.UserID != id {
			WriteErrorJSON(w, http.StatusForbidden, errors.New("forbidden"), "you can only update your own profile")
			return
		}
	}

	var req employee.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	ctx := r.Context()

	err := h.usecase.UpdateProfile(ctx, id, req)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to update employee")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "employee updated successfully")
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Assume we get the ID from the URL path, e.g., /employees/{id}
	id := r.URL.Path[len("/employees/"):]
	if id == "" {
		WriteErrorJSON(w, http.StatusBadRequest, nil, "employee ID is required")
		return
	}

	ctx := r.Context()

	err := h.usecase.Delete(ctx, id)
	if err != nil {
		if err.Error() == usecase.EmployeeNotFoundError {
			WriteErrorJSON(w, http.StatusNotFound, err, "employee not found")
			return
		}
		WriteErrorJSON(w, http.StatusInternalServerError, err, "failed to delete employee")
		return
	}

	WriteJSON(w, http.StatusOK, nil, "employee deleted successfully")
}
