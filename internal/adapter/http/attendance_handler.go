package adapterhttp

import (
	"encoding/json"
	"net/http"

	"github.com/zuyatna/shop-retail-employee-service/internal/dto/attendance"
	"github.com/zuyatna/shop-retail-employee-service/internal/usecase"
	"github.com/zuyatna/shop-retail-employee-service/internal/util/jwtutil"
)

type AttendanceHandler struct {
	attendanceUsecase *usecase.AttendanceUsecase
}

func NewAttendanceHandler(uc *usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceUsecase: uc,
	}
}

func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(UserClaimsKey).(*jwtutil.Claims)
	if !ok || claims == nil {
		WriteErrorJSON(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	var req attendance.CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteErrorJSON(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	id, err := h.attendanceUsecase.CheckIn(r.Context(), claims.UserID, req)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"id": id}, "check-in successful")
}

func (h *AttendanceHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(UserClaimsKey).(*jwtutil.Claims)
	if !ok || claims == nil {
		WriteErrorJSON(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	err := h.attendanceUsecase.CheckOut(r.Context(), claims.UserID)
	if err != nil {
		WriteErrorJSON(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "check-out successful"}, "check-out successful")
}
