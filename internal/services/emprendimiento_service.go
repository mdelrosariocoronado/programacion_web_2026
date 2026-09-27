package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/repositories"
)

var (
	ErrNombreRequerido = errors.New("el nombre del emprendimiento es obligatorio")
	ErrRubroRequerido  = errors.New("el rubro del emprendimiento es obligatorio")
	ErrNotFound        = errors.New("emprendimiento no encontrado")
)

type EmprendimientoService struct {
	repo *repositories.EmprendimientoRepository
}

func NewEmprendimientoService(repo *repositories.EmprendimientoRepository) *EmprendimientoService {
	return &EmprendimientoService{
		repo: repo,
	}
}

// DTOs de entrada y salida para desacoplar la API del driver de base de datos
type CreateEmprendimientoInput struct {
	Nombre      string `json:"nombre"`
	Url         string `json:"url"`
	Rubro       string `json:"rubro"`
	Descripcion string `json:"descripcion"`
	Logo        string `json:"logo"`
	Contacto    string `json:"contacto"`
	Activo      *bool  `json:"activo"`
}

type UpdateEmprendimientoInput struct {
	Nombre      string `json:"nombre"`
	Url         string `json:"url"`
	Rubro       string `json:"rubro"`
	Descripcion string `json:"descripcion"`
	Logo        string `json:"logo"`
	Contacto    string `json:"contacto"`
	Activo      *bool  `json:"activo"`
}

type EmprendimientoResponse struct {
	IDEmprendimiento int32  `json:"id_emprendimiento"`
	Nombre           string `json:"nombre"`
	Url              string `json:"url"`
	Rubro            string `json:"rubro"`
	Descripcion      string `json:"descripcion"`
	Logo             string `json:"logo"`
	Contacto         string `json:"contacto"`
	Activo           bool   `json:"activo"`
	CreatedAt        string `json:"created_at,omitempty"`
}

func (s *EmprendimientoService) Create(ctx context.Context, input CreateEmprendimientoInput) (EmprendimientoResponse, error) {
	input.Nombre = strings.TrimSpace(input.Nombre)
	input.Rubro = strings.TrimSpace(input.Rubro)

	if input.Nombre == "" {
		return EmprendimientoResponse{}, ErrNombreRequerido
	}
	if input.Rubro == "" {
		return EmprendimientoResponse{}, ErrRubroRequerido
	}

	activoVal := true
	if input.Activo != nil {
		activoVal = *input.Activo
	}

	params := sqlc.CreateEmprendimientoParams{
		Nombre:      input.Nombre,
		Url:         toNullString(input.Url),
		Rubro:       input.Rubro,
		Descripcion: toNullString(input.Descripcion),
		Logo:        toNullString(input.Logo),
		Contacto:    toNullString(input.Contacto),
		Activo:      sql.NullBool{Bool: activoVal, Valid: true},
	}

	creado, err := s.repo.Create(ctx, params)
	if err != nil {
		return EmprendimientoResponse{}, err
	}

	return toEmprendimientoResponse(creado), nil
}

func (s *EmprendimientoService) GetByID(ctx context.Context, id int32) (EmprendimientoResponse, error) {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EmprendimientoResponse{}, ErrNotFound
		}
		return EmprendimientoResponse{}, err
	}

	return toEmprendimientoResponse(emp), nil
}

func (s *EmprendimientoService) List(ctx context.Context) ([]EmprendimientoResponse, error) {
	lista, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]EmprendimientoResponse, 0, len(lista))
	for _, item := range lista {
		res = append(res, toEmprendimientoResponse(item))
	}

	return res, nil
}

func (s *EmprendimientoService) Update(ctx context.Context, id int32, input UpdateEmprendimientoInput) (EmprendimientoResponse, error) {
	// Verificar existencia previa
	existente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EmprendimientoResponse{}, ErrNotFound
		}
		return EmprendimientoResponse{}, err
	}

	input.Nombre = strings.TrimSpace(input.Nombre)
	input.Rubro = strings.TrimSpace(input.Rubro)

	if input.Nombre == "" {
		return EmprendimientoResponse{}, ErrNombreRequerido
	}
	if input.Rubro == "" {
		return EmprendimientoResponse{}, ErrRubroRequerido
	}

	activoVal := true
	if input.Activo != nil {
		activoVal = *input.Activo
	} else if existente.Activo.Valid {
		activoVal = existente.Activo.Bool
	}

	params := sqlc.UpdateEmprendimientoParams{
		IDEmprendimiento: id,
		Nombre:           input.Nombre,
		Url:              toNullString(input.Url),
		Rubro:            input.Rubro,
		Descripcion:      toNullString(input.Descripcion),
		Logo:             toNullString(input.Logo),
		Contacto:         toNullString(input.Contacto),
		Activo:           sql.NullBool{Bool: activoVal, Valid: true},
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return EmprendimientoResponse{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *EmprendimientoService) Delete(ctx context.Context, id int32) error {
	// Verificar existencia previa para poder devolver 404 si no existe
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// Helpers internos de mapeo
func toNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func toEmprendimientoResponse(e sqlc.Emprendimiento) EmprendimientoResponse {
	var url, desc, logo, contacto, createdAt string
	if e.Url.Valid {
		url = e.Url.String
	}
	if e.Descripcion.Valid {
		desc = e.Descripcion.String
	}
	if e.Logo.Valid {
		logo = e.Logo.String
	}
	if e.Contacto.Valid {
		contacto = e.Contacto.String
	}
	if e.CreatedAt.Valid {
		createdAt = e.CreatedAt.Time.Format(time.RFC3339)
	}

	activo := true
	if e.Activo.Valid {
		activo = e.Activo.Bool
	}

	return EmprendimientoResponse{
		IDEmprendimiento: e.IDEmprendimiento,
		Nombre:           e.Nombre,
		Url:              url,
		Rubro:            e.Rubro,
		Descripcion:      desc,
		Logo:             logo,
		Contacto:         contacto,
		Activo:           activo,
		CreatedAt:        createdAt,
	}
}
