package repositories

//consulta a la base de datos

// componentes deben depender de interfaces

import (
	"context"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

//interface para definir el contrato que UserRepository debe cumplir en su comportamiento
type UserRepo interface {

		Create(ctx context.Context, params db.CreateUsuarioParams) (db.Usuario, error)
		GetByID(ctx context.Context, id int32) (db.Usuario, error)
		ListAll(ctx context.Context) ([]db.Usuario, error)
		Update(ctx context.Context, params db.UpdateUsuarioParams) error
		Delete(ctx context.Context, id int32) error

}


type userRepository struct {
	queries *db.Queries
}

// Inyección de Dependencias
func NewUserRepository(q *db.Queries) UserRepo {
	return &userRepository{
		queries: q,
	}
}

// funciones

func (r *userRepository) Create(ctx context.Context, params db.CreateUsuarioParams) (db.Usuario, error) {
	return r.queries.CreateUsuario(ctx, params)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (db.Usuario, error) {
	return r.queries.GetUsuario(ctx, id)
}

func (r *userRepository) ListAll(ctx context.Context) ([]db.Usuario, error) {
	return r.queries.ListUsuarios(ctx)
}

func (r *userRepository) Update(ctx context.Context, params db.UpdateUsuarioParams) error {
	return r.queries.UpdateUsuario(ctx, params)
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUsuario(ctx, id)
}



//buscar por ID


//buscar por email

//listar todos


//actualizar y eliminar