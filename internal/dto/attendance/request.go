package attendance

type CheckInRequest struct {
	Location string `json:"location" binding:"required"`
}
