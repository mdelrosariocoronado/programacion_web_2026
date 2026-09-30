package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	db "emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/repositories"
)

// Errores de negocio para publicaciones
var (
	ErrTituloRequerido          = errors.New("el titulo de la publicacion es obligatorio")
	ErrTipoInvalido             = errors.New("el tipo de publicacion debe ser 'producto', 'promocion', 'servicio' u 'otro'")
	ErrPrecioInvalido           = errors.New("el precio debe ser positivo")
	ErrIDEmprendimientoInvalido = errors.New("el ID de emprendimiento debe ser mayor a cero")
	ErrIDPublicacionInvalido    = errors.New("el ID de publicacion debe ser mayor a cero")
	ErrPublicacionNoEncontrada  = errors.New("publicacion no encontrada")
)

// DTOs de entrada
type CreatePublicacionInput struct {
	IDEmprendimiento int32    `json:"id_emprendimiento"`
	Titulo           string   `json:"titulo"`
	Contenido        string   `json:"contenido"`
	ImagenUrl        string   `json:"imagen_url"`
	Tipo             string   `json:"tipo"`
	Precio           *float64 `json:"precio"`
}

type UpdatePublicacionInput struct {
	Titulo    string   `json:"titulo"`
	Contenido string   `json:"contenido"`
	ImagenUrl string   `json:"imagen_url"`
	Tipo      string   `json:"tipo"`
	Precio    *float64 `json:"precio"`
}

// PublicacionService define el contrato de logica de negocio
type PublicacionService interface {
	Create(ctx context.Context, input CreatePublicacionInput) (db.Publicacione, error)
	GetByID(ctx context.Context, id int32) (db.Publicacione, error)
	ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error)
	Update(ctx context.Context, id int32, input UpdatePublicacionInput) error
	Delete(ctx context.Context, id int32) error
}

type publicacionService struct {
	repo repositories.PublicacionRepo
}

// NewPublicacionService crea una nueva instancia de PublicacionService
func NewPublicacionService(repo repositories.PublicacionRepo) PublicacionService {
	return &publicacionService{
		repo: repo,
	}
}

// Validaciones auxiliares
func validarTipoPublicacion(tipo string) (string, error) {
	tipoNormalizado := strings.TrimSpace(strings.ToLower(tipo))
	switch tipoNormalizado {
	case "producto", "promocion", "servicio", "otro":
		return tipoNormalizado, nil
	default:
		return "", ErrTipoInvalido
	}
}

func validarYFormatearPrecio(precio *float64) (sql.NullString, error) {
	if precio == nil {
		return sql.NullString{Valid: false}, nil
	}
	if *precio <= 0 {
		return sql.NullString{}, ErrPrecioInvalido
	}
	return sql.NullString{
		String: fmt.Sprintf("%.2f", *precio),
		Valid:  true,
	}, nil
}

func toNullStringPublicacion(s string) sql.NullString {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

// Create valida los datos e inserta la publicacion
func (s *publicacionService) Create(ctx context.Context, input CreatePublicacionInput) (db.Publicacione, error) {
	if input.IDEmprendimiento <= 0 {
		return db.Publicacione{}, ErrIDEmprendimientoInvalido
	}

	titulo := strings.TrimSpace(input.Titulo)
	if titulo == "" {
		return db.Publicacione{}, ErrTituloRequerido
	}

	tipo, err := validarTipoPublicacion(input.Tipo)
	if err != nil {
		return db.Publicacione{}, err
	}

	precioNull, err := validarYFormatearPrecio(input.Precio)
	if err != nil {
		return db.Publicacione{}, err
	}

	params := db.CreatePublicacionParams{
		IDEmprendimiento: input.IDEmprendimiento,
		Titulo:           titulo,
		Contenido:        toNullStringPublicacion(input.Contenido),
		ImagenUrl:        toNullStringPublicacion(input.ImagenUrl),
		Tipo:             tipo,
		Precio:           precioNull,
	}

	pub, err := s.repo.Create(ctx, params)
	if err != nil {
		return db.Publicacione{}, fmt.Errorf("error al crear publicacion: %w", err)
	}

	return pub, nil
}

// GetByID busca una publicacion validando el ID
func (s *publicacionService) GetByID(ctx context.Context, id int32) (db.Publicacione, error) {
	if id <= 0 {
		return db.Publicacione{}, ErrIDPublicacionInvalido
	}

	pub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Publicacione{}, ErrPublicacionNoEncontrada
		}
		return db.Publicacione{}, fmt.Errorf("error al obtener publicacion: %w", err)
	}

	return pub, nil
}

// ListByEmprendimiento lista publicaciones validando el ID del emprendimiento
func (s *publicacionService) ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error) {
	if idEmprendimiento <= 0 {
		return nil, ErrIDEmprendimientoInvalido
	}

	pubs, err := s.repo.ListByEmprendimiento(ctx, idEmprendimiento)
	if err != nil {
		return nil, fmt.Errorf("error al listar publicaciones del emprendimiento: %w", err)
	}

	return pubs, nil
}

// Update valida datos y modifica la publicacion existente
func (s *publicacionService) Update(ctx context.Context, id int32, input UpdatePublicacionInput) error {
	if id <= 0 {
		return ErrIDPublicacionInvalido
	}

	// Verificar existencia previa
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPublicacionNoEncontrada
		}
		return fmt.Errorf("error al verificar existencia de publicacion: %w", err)
	}

	titulo := strings.TrimSpace(input.Titulo)
	if titulo == "" {
		return ErrTituloRequerido
	}

	tipo, err := validarTipoPublicacion(input.Tipo)
	if err != nil {
		return err
	}

	precioNull, err := validarYFormatearPrecio(input.Precio)
	if err != nil {
		return err
	}

	params := db.UpdatePublicacionParams{
		IDPublicacion: id,
		Titulo:        titulo,
		Contenido:     toNullStringPublicacion(input.Contenido),
		ImagenUrl:     toNullStringPublicacion(input.ImagenUrl),
		Tipo:          tipo,
		Precio:        precioNull,
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return fmt.Errorf("error al actualizar publicacion: %w", err)
	}

	return nil
}

// Delete elimina una publicacion validando ID y existencia
func (s *publicacionService) Delete(ctx context.Context, id int32) error {
	if id <= 0 {
		return ErrIDPublicacionInvalido
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPublicacionNoEncontrada
		}
		return fmt.Errorf("error al verificar publicacion a eliminar: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error al eliminar publicacion: %w", err)
	}

	return nil
}
