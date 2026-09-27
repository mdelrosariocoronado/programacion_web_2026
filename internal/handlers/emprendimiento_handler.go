package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"emprendimientos.com/servidor-go/internal/services"
)

type EmprendimientoHandler struct {
	service *services.EmprendimientoService
}

func NewEmprendimientoHandler(service *services.EmprendimientoService) *EmprendimientoHandler {
	return &EmprendimientoHandler{
		service: service,
	}
}

// Create: POST /api/emprendimientos
func (h *EmprendimientoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input services.CreateEmprendimientoInput
	if err := ParseJSON(r, &input); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	creado, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) || errors.Is(err, services.ErrRubroRequerido) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "Error al crear el emprendimiento: "+err.Error())
		return
	}

	SendJSON(w, http.StatusCreated, creado)
}

// GetAll: GET /api/emprendimientos
func (h *EmprendimientoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	lista, err := h.service.List(r.Context())
	if err != nil {
		SendError(w, http.StatusInternalServerError, "Error al obtener la lista de emprendimientos: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, lista)
}

// GetByID: GET /api/emprendimientos/{id}
func (h *EmprendimientoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "El parámetro 'id' debe ser un número entero positivo")
		return
	}

	emp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			SendError(w, http.StatusNotFound, "Emprendimiento no encontrado")
			return
		}
		SendError(w, http.StatusInternalServerError, "Error al consultar el emprendimiento: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, emp)
}

// Update: PUT /api/emprendimientos/{id}
func (h *EmprendimientoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "El parámetro 'id' debe ser un número entero positivo")
		return
	}

	var input services.UpdateEmprendimientoInput
	if err := ParseJSON(r, &input); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	actualizado, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			SendError(w, http.StatusNotFound, "Emprendimiento no encontrado")
			return
		}
		if errors.Is(err, services.ErrNombreRequerido) || errors.Is(err, services.ErrRubroRequerido) {
			SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendError(w, http.StatusInternalServerError, "Error al actualizar el emprendimiento: "+err.Error())
		return
	}

	SendJSON(w, http.StatusOK, actualizado)
}

// Delete: DELETE /api/emprendimientos/{id}
func (h *EmprendimientoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "El parámetro 'id' debe ser un número entero positivo")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			SendError(w, http.StatusNotFound, "Emprendimiento no encontrado")
			return
		}
		SendError(w, http.StatusInternalServerError, "Error al eliminar el emprendimiento: "+err.Error())
		return
	}

	// 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// Helper para extraer y convertir el parámetro {id} de la URL
func parseID(r *http.Request) (int32, error) {
	idStr := r.PathValue("id")
	val, err := strconv.Atoi(idStr)
	if err != nil || val <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(val), nil
}
