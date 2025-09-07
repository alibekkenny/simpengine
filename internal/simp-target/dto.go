package simptarget

type SimpTargetRequestDTO struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"required"`
}

type CreateSimpTargetResponseDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
