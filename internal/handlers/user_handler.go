package handlers


import (
	"net/http"

	"emprendimientos.com/servidor-go/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

// Nota la N mayúscula
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Página de inicio"))
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint de usuarios"))
}