package handlers

import (
	"net/http"

	"emprendimientos.com/servidor-go/views"
)

// HomeHandler maneja la página de inicio del portal.
type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// RenderPage: GET / (SSR con Templ - página de inicio)
func (h *HomeHandler) RenderPage(w http.ResponseWriter, r *http.Request) {
	if err := views.HomePage().Render(r.Context(), w); err != nil {
		http.Error(w, "Error al renderizar la página de inicio: "+err.Error(), http.StatusInternalServerError)
	}
}
