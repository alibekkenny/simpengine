package user

type RegisterRequestDTO struct {
	Login    string `json:"login" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
	Token    string `json:"token"` //optional (required only for admins)
}

type RegisterResponseDTO struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
