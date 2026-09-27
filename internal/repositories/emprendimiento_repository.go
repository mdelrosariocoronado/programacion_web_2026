package repositories

import (
	"context"

	"emprendimientos.com/servidor-go/db/sqlc"
)

type EmprendimientoRepository struct {
	queries *sqlc.Queries
}

func NewEmprendimientoRepository(queries *sqlc.Queries) *EmprendimientoRepository {
	return &EmprendimientoRepository{
		queries: queries,
	}
}

func (r *EmprendimientoRepository) Create(ctx context.Context, params sqlc.CreateEmprendimientoParams) (sqlc.Emprendimiento, error) {
	return r.queries.CreateEmprendimiento(ctx, params)
}

func (r *EmprendimientoRepository) GetByID(ctx context.Context, id int32) (sqlc.Emprendimiento, error) {
	return r.queries.GetEmprendimiento(ctx, id)
}

func (r *EmprendimientoRepository) List(ctx context.Context) ([]sqlc.Emprendimiento, error) {
	return r.queries.ListEmprendimientos(ctx)
}

func (r *EmprendimientoRepository) Update(ctx context.Context, params sqlc.UpdateEmprendimientoParams) error {
	return r.queries.UpdateEmprendimiento(ctx, params)
}

func (r *EmprendimientoRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteEmprendimiento(ctx, id)
}
