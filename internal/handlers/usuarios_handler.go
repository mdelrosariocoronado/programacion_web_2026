package handlers

// VALIDA HTTP, decodifica JSON, llama al service de usuario, escribe JSON y status code



//lee JSON r.BOdy

//MANEJO DE PATH VAALUE

//invocar service

//mapear resultado http


import (
	"net/http"
	"encoding/json"
	"errors"

	"emprendimientos.com/servidor-go/internal/services"
)

type UserHandler struct {
	service services.UserService
}

// DTOs para no acoplar la estructura HTTP a la estructura de la base de datos
type CreateUserRequest struct {
	NombreCompleto   string `json:"nombre_completo"`
	Email            string `json:"email"`
	Clave            string `json:"clave"`
	Rol              string `json:"rol"`
	IdEmprendimiento *int32 `json:"id_emprendimiento,omitempty"`
}

type UpdateUserRequest struct {
	NombreCompleto   string `json:"nombre_completo"`
	Email            string `json:"email"`
	Clave            string `json:"clave"`
	Rol              string `json:"rol"`
	IdEmprendimiento *int32 `json:"id_emprendimiento,omitempty"`
}





// para inyeccion de dependencias
func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

// POST el api para USUARIOS


func (h *UserHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Página de inicio"))
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint de usuarios"))
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "El cuerpo del JSON no es valido", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.CrearUsuario(
		r.Context(),
		req.NombreCompleto,
		req.Rol,
		req.Email,
		req.Clave,
		req.IdEmprendimiento,
	)
	if err != nil{
		if errors.Is(err, services.ErrFormatoInvalidoEmail) || errors.Is(err, services.ErrClave) || errors.Is(err, services.ErrRolInvalido) || errors.Is(err, services.ErrNombreVacio) || errors.Is(err, services.ErrEmailDuplicado){
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // mensaje 201
	json.NewEncoder(w).Encode(usuario)
	w.Write([]byte(`{"mensaje":"Usuario creado"}`))
}


