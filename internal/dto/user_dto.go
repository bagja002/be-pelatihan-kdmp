package dto

type CreateUserRequest struct {
	Nama     string `json:"nama" validate:"required,max=255"`
	Username string `json:"username" validate:"required,max=128"`
	Password string `json:"password" validate:"required,min=6,max=128"`
	Type     string `json:"type" validate:"required,oneof=super_admin admin"`
	IDSatdik *uint  `json:"idSatdik" validate:"omitempty"`
}

type UpdateUserRequest struct {
	Nama     string `json:"nama" validate:"omitempty,max=255"`
	Password string `json:"password" validate:"omitempty,min=6,max=128"`
	Type     string `json:"type" validate:"omitempty,oneof=super_admin admin"`
	IDSatdik *uint  `json:"idSatdik" validate:"omitempty"`
}
