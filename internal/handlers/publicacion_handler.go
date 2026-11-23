package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/services"
	"emprendimientos.com/servidor-go/views"
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

// RenderPage: GET /publicaciones (SSR con Templ, soporta filtro ?id_emprendimiento=)
func (h *PublicacionHandler) RenderPage(w http.ResponseWriter, r *http.Request) {
	h.ListarHTML(w, r)
}

// CreateFromForm: POST /publicaciones
func (h *PublicacionHandler) CreateFromForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendHTMXError(w, http.StatusBadRequest, "Error al procesar los datos del formulario")
		return
	}

	idEmpStr := strings.TrimSpace(r.FormValue("id_emprendimiento"))
	idEmp, err := strconv.Atoi(idEmpStr)
	if err != nil || idEmp <= 0 {
		SendHTMXError(w, http.StatusBadRequest, "El ID de emprendimiento debe ser un número entero positivo")
		return
	}

	input := services.CreatePublicacionInput{
		IDEmprendimiento: int32(idEmp),
		Titulo:           r.FormValue("titulo"),
		Tipo:             r.FormValue("tipo"),
		Contenido:        r.FormValue("contenido"),
		ImagenUrl:        r.FormValue("imagen_url"),
	}

	if precioStr := strings.TrimSpace(r.FormValue("precio")); precioStr != "" {
		precio, err := strconv.ParseFloat(precioStr, 64)
		if err != nil {
			SendHTMXError(w, http.StatusBadRequest, "El precio debe ser un número válido")
			return
		}
		input.Precio = &precio
	}

	_, err = h.service.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrTituloRequerido) ||
			errors.Is(err, services.ErrTipoInvalido) ||
			errors.Is(err, services.ErrPrecioInvalido) ||
			errors.Is(err, services.ErrIDEmprendimientoInvalido) {
			SendHTMXError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al crear la publicación: "+err.Error())
		return
	}

	publicaciones, err := h.service.List(r.Context())
	if err != nil {
		SendHTMXError(w, http.StatusInternalServerError, "Error al obtener la lista actualizada: "+err.Error())
		return
	}
	if err := views.PublicacionList(publicaciones).Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la lista: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := views.PublicacionFormOOB().Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar el formulario limpio: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := views.ErrorBannerOOB().Render(r.Context(), w); err != nil {
		http.Error(w, "Error al limpiar alertas: "+err.Error(), http.StatusInternalServerError)
	}
}

// DeleteFromHTMX: DELETE /publicaciones/{id} (Borrado reactivo sin recarga de página)
func (h *PublicacionHandler) DeleteFromHTMX(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		SendHTMXError(w, http.StatusBadRequest, "El parámetro 'id' debe ser un número entero positivo")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, services.ErrPublicacionNoEncontrada) {
			SendHTMXError(w, http.StatusNotFound, "Publicación no encontrada")
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al eliminar la publicación: "+err.Error())
		return
	}

	// HTMX: 200 OK con cuerpo vacío reemplaza el target (closest tr) por nada, eliminando la fila del DOM
	w.WriteHeader(http.StatusOK)
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
func (h *PublicacionHandler) ListarHTML(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var (
		publicaciones []sqlc.Publicacione
		err           error
	)

	param := r.URL.Query().Get("id_emprendimiento")
	id, errConv := strconv.Atoi(param)

	// Si viene el parámetro en la URL y es un número mayor a 0, filtra
	if param != "" && errConv == nil && id > 0 {
		publicaciones, err = h.service.ListByEmprendimiento(ctx, int32(id))
	} else {
		// De lo contrario, trae todas las publicaciones
		publicaciones, err = h.service.List(ctx)
	}

	if err != nil {
		http.Error(w, "Error al cargar publicaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	publicaciones = filtrarPublicaciones(publicaciones, query)

	if isHTMXPartial(r) {
		if err := views.PublicacionList(publicaciones).Render(ctx, w); err != nil {
			http.Error(w, "Error al renderizar la lista: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := views.PublicacionesPage(publicaciones, query).Render(ctx, w); err != nil {
		http.Error(w, "Error al renderizar: "+err.Error(), http.StatusInternalServerError)
	}
}

func filtrarPublicaciones(pubs []sqlc.Publicacione, query string) []sqlc.Publicacione {
	if query == "" {
		return pubs
	}
	q := strings.ToLower(query)
	res := make([]sqlc.Publicacione, 0, len(pubs))
	for _, p := range pubs {
		contenido := ""
		if p.Contenido.Valid {
			contenido = p.Contenido.String
		}
		if strings.Contains(strings.ToLower(p.Titulo), q) ||
			strings.Contains(strings.ToLower(p.Tipo), q) ||
			strings.Contains(strings.ToLower(contenido), q) {
			res = append(res, p)
		}
	}
	return res
}