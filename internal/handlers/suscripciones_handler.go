package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/services"
	"emprendimientos.com/servidor-go/views"
)

type emprendimientoLister interface {
	List(ctx context.Context) ([]services.EmprendimientoResponse, error)
}

// SuscripcionHandler maneja las peticiones HTTP de la entidad suscripciones
type SuscripcionHandler struct {
	service     services.SuscripcionService
	userService services.UserService
	empService  emprendimientoLister
}

func NewSuscripcionHandler(s services.SuscripcionService, u services.UserService, e emprendimientoLister) *SuscripcionHandler {
	return &SuscripcionHandler{
		service:     s,
		userService: u,
		empService:  e,
	}
}

// RenderPage: GET /suscripciones
func (h *SuscripcionHandler) RenderPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usuarios, err := h.userService.ListarUsuarios(ctx)
	if err != nil {
		if isHTMXPartial(r) {
			SendHTMXError(w, http.StatusInternalServerError, "Error al obtener usuarios: "+err.Error())
		} else {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	selectedUserID := resolverUsuarioSeleccionado(r.URL.Query().Get("id_usuario"), usuarios)
	emps, err := h.empService.List(ctx)
	if err != nil {
		if isHTMXPartial(r) {
			SendHTMXError(w, http.StatusInternalServerError, "Error al obtener emprendimientos: "+err.Error())
		} else {
			http.Error(w, "Error al obtener emprendimientos: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	subsMap := make(map[int32]db.Suscripcione)
	if selectedUserID > 0 {
		subs, err := h.service.ListByUsuario(ctx, selectedUserID)
		if err != nil {
			if isHTMXPartial(r) {
				SendHTMXError(w, http.StatusInternalServerError, "Error al obtener suscripciones: "+err.Error())
			} else {
				http.Error(w, "Error al obtener suscripciones: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		for _, s := range subs {
			subsMap[s.IDEmprendimiento] = s
		}
	}

	if isHTMXPartial(r) {
		if err := views.SuscripcionesPanel(selectedUserID, emps, subsMap).Render(ctx, w); err != nil {
			http.Error(w, "Error al renderizar suscripciones: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := views.SuscripcionesPage(usuarios, selectedUserID, emps, subsMap).Render(ctx, w); err != nil {
		http.Error(w, "Error al renderizar la vista: "+err.Error(), http.StatusInternalServerError)
	}
}

// SubscribeFromHTMX: POST /suscripciones/{id_usuario}/{id_emprendimiento}
func (h *SuscripcionHandler) SubscribeFromHTMX(w http.ResponseWriter, r *http.Request) {
	idUsuario, idEmprendimiento, err := parseSuscripcionPathIDs(r)
	if err != nil {
		SendHTMXError(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := h.service.Suscribe(r.Context(), idUsuario, idEmprendimiento); err != nil {
		if errors.Is(err, services.ErrYaSuscrito) {
			SendHTMXError(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, services.ErrDatosSuscripcionInvalidos) {
			SendHTMXError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al registrar suscripción: "+err.Error())
		return
	}

	h.renderPanelActualizado(w, r, idUsuario)
}

// CancelFromHTMX: DELETE /suscripciones/{id_usuario}/{id_emprendimiento}
func (h *SuscripcionHandler) CancelFromHTMX(w http.ResponseWriter, r *http.Request) {
	idUsuario, idEmprendimiento, err := parseSuscripcionPathIDs(r)
	if err != nil {
		SendHTMXError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Cancel(r.Context(), idUsuario, idEmprendimiento); err != nil {
		if errors.Is(err, services.ErrSuscripcionNoEncontrada) {
			SendHTMXError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, services.ErrDatosSuscripcionInvalidos) {
			SendHTMXError(w, http.StatusBadRequest, err.Error())
			return
		}
		SendHTMXError(w, http.StatusInternalServerError, "Error al cancelar suscripción: "+err.Error())
		return
	}

	h.renderPanelActualizado(w, r, idUsuario)
}

func (h *SuscripcionHandler) renderPanelActualizado(w http.ResponseWriter, r *http.Request, idUsuario int32) {
	ctx := r.Context()
	emps, err := h.empService.List(ctx)
	if err != nil {
		SendHTMXError(w, http.StatusInternalServerError, "Error al cargar emprendimientos: "+err.Error())
		return
	}

	subs, err := h.service.ListByUsuario(ctx, idUsuario)
	if err != nil {
		SendHTMXError(w, http.StatusInternalServerError, "Error al actualizar suscripciones: "+err.Error())
		return
	}

	subsMap := make(map[int32]db.Suscripcione, len(subs))
	for _, s := range subs {
		subsMap[s.IDEmprendimiento] = s
	}

	if err := views.SuscripcionesPanel(idUsuario, emps, subsMap).Render(ctx, w); err != nil {
		http.Error(w, "Error al renderizar suscripciones: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := views.ErrorBannerOOB().Render(ctx, w); err != nil {
		http.Error(w, "Error al limpiar alertas: "+err.Error(), http.StatusInternalServerError)
	}
}

func resolverUsuarioSeleccionado(param string, usuarios []db.Usuario) int32 {
	if len(usuarios) == 0 {
		return 0
	}
	if id, err := strconv.Atoi(strings.TrimSpace(param)); err == nil && id > 0 {
		for _, u := range usuarios {
			if u.IDUsuario == int32(id) {
				return u.IDUsuario
			}
		}
	}
	return usuarios[0].IDUsuario
}

func parseSuscripcionPathIDs(r *http.Request) (int32, int32, error) {
	idUser, err := strconv.Atoi(r.PathValue("id_usuario"))
	if err != nil || idUser <= 0 {
		return 0, 0, errors.New("el parámetro 'id_usuario' debe ser un número entero positivo")
	}
	idEmp, err := strconv.Atoi(r.PathValue("id_emprendimiento"))
	if err != nil || idEmp <= 0 {
		return 0, 0, errors.New("el parámetro 'id_emprendimiento' debe ser un número entero positivo")
	}
	return int32(idUser), int32(idEmp), nil
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
