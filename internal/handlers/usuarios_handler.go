package handlers

// VALIDA HTTP, decodifica JSON, llama al service de usuario, escribe JSON y status code

//lee JSON r.BOdy

//MANEJO DE PATH VAALUE

//invocar service

//mapear resultado http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		req.Email,
		req.Clave,
		req.Rol,
		req.IdEmprendimiento,
	)
	if err != nil{
		if errors.Is(err, services.ErrFormatoInvalidoEmail) || errors.Is(err, services.ErrClave) || errors.Is(err, services.ErrRolInvalido) || errors.Is(err, services.ErrNombreVacio) || errors.Is(err, services.ErrEmailDuplicado){
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "error interno del servidor", http.StatusInternalServerError)
		return
		
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // mensaje 201
	json.NewEncoder(w).Encode(usuario)
}

// metodo con GET en api de usuarios cin id

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request){

	//con go 1.22 r.pathvaues extrae partes del path
	id_usuario := r.PathValue("id")
	id, err := strconv.ParseInt(id_usuario, 10, 32)
	if err != nil {
		http.Error(w, " ID de usuario invalido", http.StatusBadRequest)
		return
	}

	usuario, err := h.service.ObtenerUsuario(r.Context(), int32(id))

	if err!= nil {
		//usuario que no esta en el registro, se tira error especial por  sqlc
		http.Error(w, "El usuario ingresado no esta registrado", http.StatusNotFound) // error 404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200, se acepto
	json.NewEncoder(w).Encode(usuario)

}

// listar los usuarios
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.service.ListarUsuarios(r.Context())
	if err != nil {
		http.Error(w, "HUbo un error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 esta bien
	json.NewEncoder(w).Encode(usuarios)
}


// PUT /api/usuarios/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") //para parsear rutas  como /api/usuarios/{id}
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Identificador de usuario inválido", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Cuerpo de solicitud JSON inválido", http.StatusBadRequest)
		return
	}

	err = h.service.ActualizarUsuario(
		r.Context(),
		int32(id),
		req.NombreCompleto,
		req.Email,
		req.Clave,
		req.Rol,
		req.IdEmprendimiento,
	)
	if err != nil {
		if errors.Is(err, services.ErrClave) ||
			errors.Is(err, services.ErrEmailDuplicado) ||
			errors.Is(err, services.ErrFormatoInvalidoEmail) ||
			errors.Is(err, services.ErrNombreVacio) ||
			errors.Is(err, services.ErrRolInvalido) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Usuario actualizado exitosamente"})
}

// DELETE /api/usuarios/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Identificador de usuario inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.EliminarUsuario(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	// 204 No Content no lleva body en la respuesta
	w.WriteHeader(http.StatusNoContent)
}