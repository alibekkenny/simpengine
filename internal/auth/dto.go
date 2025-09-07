package auth

type LoginRequestDTO struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}
