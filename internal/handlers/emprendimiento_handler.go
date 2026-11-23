package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"emprendimientos.com/servidor-go/internal/services"
	"emprendimientos.com/servidor-go/views"
)

type EmprendimientoHandler struct {
	service *services.EmprendimientoService
}

func NewEmprendimientoHandler(service *services.EmprendimientoService) *EmprendimientoHandler {
	return &EmprendimientoHandler{
		service: service,
	}
}

// RenderPage: GET /emprendimientos
func (h *EmprendimientoHandler) RenderPage(w http.ResponseWriter, r *http.Request) {
	lista, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener la lista de emprendimientos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	lista = filtrarEmprendimientos(lista, query)

	if isHTMXPartial(r) {
		if err := views.EmprendimientoList(lista).Render(r.Context(), w); err != nil {
			http.Error(w, "Error al renderizar la lista: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := views.EmprendimientosPage(lista, query).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la vista: "+err.Error(), http.StatusInternalServerError)
	}
}

// CreateFromForm: POST /emprendimientos
func (h *EmprendimientoHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendHTMXError(w, http.StatusBadRequest, "Error al procesar los datos del formulario")
		return
	}

	input := services.CreateEmprendimientoInput{
		Nombre:      r.FormValue("nombre"),
		Rubro:       r.FormValue("rubro"),
		Descripcion: r.FormValue("descripcion"),
		Contacto:    r.FormValue("contacto"),
		Url:         r.FormValue("url"),
		Logo:        r.FormValue("logo"),
	}

	_, err := h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) || errors.Is(err, services.ErrRubroRequerido) {
			SendHTMXError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al crear el emprendimiento: "+err.Error())
		return
	}

	lista, err := h.service.List(r.Context())
	if err != nil {
		SendHTMXError(w, http.StatusInternalServerError, "Error al obtener la lista actualizada: "+err.Error())
		return
	}
	if err := views.EmprendimientoList(lista).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la lista: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := views.EmprendimientoFormOOB().Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar el formulario limpio: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := views.ErrorBannerOOB().Render(r.Context(), w); err != nil {
		http.Error(w, "Error al limpiar alertas: "+err.Error(), http.StatusInternalServerError)
	}
}

// DeleteFromHTMX: DELETE /emprendimientos/{id} (Borrado reactivo sin recarga de página)
func (h *EmprendimientoHandler) DeleteFromHTMX(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendHTMXError(w, http.StatusBadRequest, "El parámetro 'id' debe ser un número entero positivo")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			SendHTMXError(w, http.StatusNotFound, "Emprendimiento no encontrado")
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al eliminar el emprendimiento: "+err.Error())
		return
	}

	// HTMX: 200 OK con cuerpo vacío reemplaza el target (closest tr) por nada, eliminando la fila del DOM
	w.WriteHeader(http.StatusOK)
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

func filtrarEmprendimientos(items []services.EmprendimientoResponse, query string) []services.EmprendimientoResponse {
	if query == "" {
		return items
	}
	q := strings.ToLower(query)
	res := make([]services.EmprendimientoResponse, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Nombre), q) ||
			strings.Contains(strings.ToLower(item.Rubro), q) ||
			strings.Contains(strings.ToLower(item.Descripcion), q) ||
			strings.Contains(strings.ToLower(item.Contacto), q) {
			res = append(res, item)
		}
	}
	return res
}
