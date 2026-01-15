package employee

type CreateEmployeeRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Role        string `json:"role" validate:"required,oneof=admin staff manager"`
	Position    string `json:"position" validate:"required"`
	Salary      int    `json:"salary" validate:"required,gt=0"`
	Status      string `json:"status" validate:"required,oneof=active inactive"`
	BirthDate   string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Address     string `json:"address" validate:"required"`
	City        string `json:"city" validate:"required"`
	Province    string `json:"province" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}
