package employee

type UpdateEmployeeRequest struct {
	Name        *string `json:"name,omitempty"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	Password    *string `json:"password,omitempty" validate:"omitempty,min=8"`
	Role        *string `json:"role,omitempty" validate:"omitempty,oneof=admin staff manager"`
	Position    *string `json:"position,omitempty"`
	Salary      *int    `json:"salary,omitempty" validate:"omitempty,gt=0"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	BirthDate   *string `json:"birth_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	Province    *string `json:"province,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	Photo       *string `json:"photo,omitempty"`
}
