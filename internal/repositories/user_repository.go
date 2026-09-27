package repositories

import (
	"emprendimientos.com/servidor-go/db/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

// Nota la N mayúscula para exportar la función
func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}