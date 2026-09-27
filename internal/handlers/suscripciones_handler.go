package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"emprendimientos.com/servidor-go/internal/services"
)

// SuscripcionHandler maneja las peticiones HTTP de la entidad suscripciones
type SuscripcionHandler struct {
	service services.SuscripcionService
}

func NewSuscripcionHandler(s services.SuscripcionService) *SuscripcionHandler {
	return &SuscripcionHandler{service: s}
}

// DTO para la creación y baja de suscripciones
type SuscripcionRequestDTO struct {
	IDUsuario        int32 `json:"id_usuario"`
	IDEmprendimiento int32 `json:"id_emprendimiento"`
}

// POST /api/suscripciones
func (h *SuscripcionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto SuscripcionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	sub, err := h.service.Suscribe(r.Context(), dto.IDUsuario, dto.IDEmprendimiento)
	if err != nil {
		if errors.Is(err, services.ErrYaSuscrito) {
			http.Error(w, err.Error(), http.StatusConflict) // 409
			return
		}
		if errors.Is(err, services.ErrDatosSuscripcionInvalidos) {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// GET /api/usuarios/{id}/suscripciones
func (h *SuscripcionHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || id <= 0 {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	subs, err := h.service.ListByUsuario(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al obtener suscripciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subs)
}

// DELETE /api/suscripciones
func (h *SuscripcionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var dto SuscripcionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.Cancel(r.Context(), dto.IDUsuario, dto.IDEmprendimiento); err != nil {
		if errors.Is(err, services.ErrSuscripcionNoEncontrada) {
			http.Error(w, err.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}