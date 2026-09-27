package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/repositories"
)

// variables para tener los errores en el servicio 
var (
	ErrYaSuscrito                = errors.New("el usuario ya esta suscrito al emprendimiento")
	ErrSuscripcionNoEncontrada   = errors.New("la suscripción no fue encontrada")
	ErrDatosSuscripcionInvalidos = errors.New("los identificadores de usuario y emprendimiento deben ser mayores a cero")
)

// SuscripcionService tiene logica de negocio
type SuscripcionService interface {
	Suscribe(ctx context.Context, idUsuario, idEmprendimiento int32) (db.Suscripcione, error)
	ListByUsuario(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error)
	Cancel(ctx context.Context, idUsuario, idEmprendimiento int32) error
}

type suscripcionService struct {
	repo repositories.SuscripcionRepo
}

//  inyección de dependencias
func NewSuscripcionService(repo repositories.SuscripcionRepo) SuscripcionService {
	return &suscripcionService{
		repo: repo,
	}
}

// REGLA unicidad y PERSISTENCIA DEL N:M
func (s *suscripcionService) Suscribe(ctx context.Context, idUsuario, idEmprendimiento int32) (db.Suscripcione, error) {
	// identificadores validos
	if idUsuario <= 0 || idEmprendimiento <= 0 {
		return db.Suscripcione{}, ErrDatosSuscripcionInvalidos
	}

	// verificar si existe suscripcion
	_, err := s.repo.GetByUserEmprendim(ctx, idUsuario, idEmprendimiento)
	if err == nil {
		return db.Suscripcione{}, ErrYaSuscrito //  existe
	} else if !errors.Is(err, sql.ErrNoRows) {
		return db.Suscripcione{}, fmt.Errorf("error al verificar existencia de suscripción: %w", err)
	}

	// parámetros para sqlc para el repo
	params := db.CreateSuscripcionParams{
		IDUsuario:        idUsuario,
		IDEmprendimiento: idEmprendimiento,
	}

	sub, err := s.repo.Create(ctx, params)
	if err != nil {
		return db.Suscripcione{}, fmt.Errorf("error al registrar suscripción: %w", err)
	}

	return sub, nil
}

// devolver todas las suscripciones de un usuario
func (s *suscripcionService) ListByUsuario(ctx context.Context, idUsuario int32) ([]db.Suscripcione, error) {
	if idUsuario <= 0 {
		return nil, ErrDatosSuscripcionInvalidos
	}

	subs, err := s.repo.List(ctx, idUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al listar suscripciones del usuario %d: %w", idUsuario, err)
	}

	return subs, nil
}

// antes de dar de baja chequear para cancelar
func (s *suscripcionService) Cancel(ctx context.Context, idUsuario, idEmprendimiento int32) error {
	if idUsuario <= 0 || idEmprendimiento <= 0 {
		return ErrDatosSuscripcionInvalidos
	}

	// Verificar si la relación existe antes de intentar eliminarla
	_, err := s.repo.GetByUserEmprendim(ctx, idUsuario, idEmprendimiento)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSuscripcionNoEncontrada
		}
		return fmt.Errorf("error al comprobar suscripción a cancelar: %w", err)
	}

	params := db.DeleteSuscripcionParams{
		IDUsuario:        idUsuario,
		IDEmprendimiento: idEmprendimiento,
	}

	if err := s.repo.Delete(ctx, params); err != nil {
		return fmt.Errorf("error al dar de baja suscripción: %w", err)
	}

	return nil
}