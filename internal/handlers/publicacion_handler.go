package handlers

import (
	"errors"
	"net/http"

	"emprendimientos.com/servidor-go/internal/services"
)

// PublicacionHandler maneja las peticiones HTTP de la entidad publicaciones
type PublicacionHandler struct {
	service services.PublicacionService
}

// NewPublicacionHandler crea una nueva instancia de PublicacionHandler
func NewPublicacionHandler(service services.PublicacionService) *PublicacionHandler {
	return &PublicacionHandler{
		service: service,
	}
}

// Create: POST /api/publicaciones (201 Created)
func (h *PublicacionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input services.CreatePublicacionInput
	if err := ParseJSON(r, &input); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	pub, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrTituloRequerido) ||
			errors.Is(err, services.ErrTipoInvalido) ||
			errors.Is(err, services.ErrPrecioInvalido) ||
			errors.Is(err, services.ErrIDEmprendimientoInvalido) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "error al crear la publicacion: "+err.Error())
		return
	}

	SendJSON(w, http.StatusCreated, pub)
}

// GetByID: GET /api/publicaciones/{id} (200 OK / 404 Not Found)
func (h *PublicacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "el parametro 'id' debe ser un numero entero positivo")
		return
	}

	pub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPublicacionNoEncontrada) {
			SendError(w, http.StatusNotFound, "publicacion no encontrada")
			return
		}
		SendError(w, http.StatusInternalServerError, "error al consultar la publicacion: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, pub)
}

// ListByEmprendimiento: GET /api/emprendimientos/{id}/publicaciones (200 OK)
func (h *PublicacionHandler) ListByEmprendimiento(w http.ResponseWriter, r *http.Request) {
	idEmprendimiento, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "el parametro 'id' del emprendimiento debe ser un numero entero positivo")
		return
	}

	pubs, err := h.service.ListByEmprendimiento(r.Context(), idEmprendimiento)
	if err != nil {
		if errors.Is(err, services.ErrIDEmprendimientoInvalido) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "error al listar las publicaciones: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, pubs)
}

// Update: PUT /api/publicaciones/{id} (200 OK)
func (h *PublicacionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "el parametro 'id' debe ser un numero entero positivo")
		return
	}

	var input services.UpdatePublicacionInput
	if err := ParseJSON(r, &input); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, services.ErrPublicacionNoEncontrada) {
			SendError(w, http.StatusNotFound, "publicacion no encontrada")
			return
		}
		if errors.Is(err, services.ErrTituloRequerido) ||
			errors.Is(err, services.ErrTipoInvalido) ||
			errors.Is(err, services.ErrPrecioInvalido) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "error al actualizar la publicacion: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, map[string]string{"mensaje": "publicacion actualizada correctamente"})
}

// Delete: DELETE /api/publicaciones/{id} (204 No Content)
func (h *PublicacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "el parametro 'id' debe ser un numero entero positivo")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPublicacionNoEncontrada) {
			SendError(w, http.StatusNotFound, "publicacion no encontrada")
			return
		}
		SendError(w, http.StatusInternalServerError, "error al eliminar la publicacion: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
