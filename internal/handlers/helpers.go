package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

// ResponseError estructura uniforme para respuestas de error en la API
type ResponseError struct {
	Error string `json:"error"`
}

// SendJSON serializa y envía una respuesta en formato JSON con el código HTTP correspondiente
func SendJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"Error interno al procesar la respuesta"}`, http.StatusInternalServerError)
		}
	}
}

// SendError responde con un objeto JSON estándar {"error": "mensaje"}
func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendJSON(w, statusCode, ResponseError{Error: message})
}

// Redirige errores HTMX al banner #htmx-error para no reemplazar la tabla ni el formulario
func SendHTMXError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Retarget", "#htmx-error")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.WriteHeader(statusCode)
	fmt.Fprint(w, html.EscapeString(message))
}

// Distingue peticiones parciales HTMX (ej. buscador) de navegación completa con hx-boost
func isHTMXPartial(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Boosted") != "true"
}

// ParseJSON decodifica el cuerpo de la petición HTTP a la estructura indicada
func ParseJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("El cuerpo de la solicitud no puede estar vacío")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // bloquea campos no definidos en la struct

	if err := decoder.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("JSON mal formado en el byte %d", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("Tipo de dato incorrecto para el campo '%s'", unmarshalTypeError.Field)
		case errors.Is(err, io.EOF):
			return errors.New("El cuerpo de la solicitud está vacío")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			campo := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("Campo no permitido en el JSON: %s", campo)
		default:
			return fmt.Errorf("Error al procesar JSON: %w", err)
		}
	}

	return nil
}
