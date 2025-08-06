package simptarget

import "net/http"

type SimpTargetHandler struct {
}

func NewSimpTargetHandler() *SimpTargetHandler {
	return &SimpTargetHandler{}
}

func (s *SimpTargetHandler) ViewSimpTarget(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
