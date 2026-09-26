package handlers

// VALIDA HTTP, decodifica JSON, llama al service de usuario, escribe JSON y status code



import (
	"net/http"

	"emprendimientos.com/servidor-go/internal/services"
)

type UserHandler struct {
	service services.UserService
}

// Nota la N mayúscula
func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Página de inicio"))
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint de usuarios"))
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"mensaje":"Usuario creado"}`))
}
