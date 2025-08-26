package simptarget

type SimpTarget struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int    `json:"-"`
}
