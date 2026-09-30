package repositories

import (
	"context"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

// PublicacionRepo define el contrato para las operaciones de publicaciones en la BD
type PublicacionRepo interface {
	Create(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error)
	GetByID(ctx context.Context, id int32) (db.Publicacione, error)
	ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error)
	Update(ctx context.Context, params db.UpdatePublicacionParams) error
	Delete(ctx context.Context, id int32) error
}

// publicacionRepository implementa PublicacionRepo interactuando con sqlc
type publicacionRepository struct {
	queries *db.Queries
}

// NewPublicacionRepository crea una nueva instancia de PublicacionRepo
func NewPublicacionRepository(q *db.Queries) PublicacionRepo {
	return &publicacionRepository{
		queries: q,
	}
}

// Create inserta una nueva publicación
func (r *publicacionRepository) Create(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error) {
	return r.queries.CreatePublicacion(ctx, params)
}

// GetByID busca una publicación por su ID
func (r *publicacionRepository) GetByID(ctx context.Context, id int32) (db.Publicacione, error) {
	return r.queries.GetPublicacion(ctx, id)
}

// ListByEmprendimiento obtiene todas las publicaciones asociadas a un emprendimiento
func (r *publicacionRepository) ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error) {
	return r.queries.ListPublicacionesByEmprendimiento(ctx, idEmprendimiento)
}

// Update modifica los datos de una publicación existente
func (r *publicacionRepository) Update(ctx context.Context, params db.UpdatePublicacionParams) error {
	return r.queries.UpdatePublicacion(ctx, params)
}

// Delete elimina una publicación por su ID
func (r *publicacionRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeletePublicacion(ctx, id)
}
