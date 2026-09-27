package repositories

import (
	"context"
	db "emprendimientos.com/servidor-go/db/sqlc"
)



//interfaz publica -> operaiones necesarias para gestionar tabla de suscripciones


type SuscripcionRepo interface {
	//CreateSuscripcion 
	Create(ctx context.Context, params db.CreateSuscripcionParams) (db.Suscripcione, error)
	//GetSuscripcionByUserAndEmprendimiento 
	GetByUserEmprendim(ctx context.Context, id_user int32, id_emp int32) (db.Suscripcione, error)
	//ListSuscripcionesByUsuario
	List(ctx context.Context) ([]db.Suscripcione, error)
	//UpdateSuscripcion
	Update(ctx context.Context) ([]db.Suscripcione, error)  //solo debe dar error, no se pueden modificar, solo se dan de alta y de baja es N:M
	//DeleteSuscripcion 
	Delete(ctx context.Context, id int32) error

}


//estructura privada que use *db.QUeries para llamar a funciones de sqlc

type suscRepository struct {
	queries *db.Queries
}

//newsuscrotion repositry ...

func NewSuscripcionRepository(q *db.Queries) SuscripcionRepo {
	return &suscRepository{queries: q}
}

func (r *suscRepository) Create(ctx context.Context, params db.CreateSuscripcionParams) (db.Suscripcione, error) {
	return r.queries.CreateSuscripcion(ctx, params)
}

func (r *suscRepository) GetByUserEmprendim(ctx context.Context, id_usuario int32, id_emp int32) (db.Suscripcione, error) {
	return r.queries.GetSuscripcionByUserAndEmprendimiento(ctx, db.GetSuscripcionByUserAndEmprendimientoParams{
		IDUsuario: id_usuario,
		IDEmprendimiento: id_emp,
	})
}

func (r *suscRepository) List(ctx context.Context, id_usuario int32) ([]db.Suscripcione, error) {
	return r.queries.ListSuscripcionesByUsuario(ctx, id_usuario)
}

func (r *suscRepository) Update(ctx context.Context, params db.UpdateSuscripcionParams) (db.Suscripcione, error) {
	return r.queries.UpdateSuscripcion(ctx, params)
}

func (r *suscRepository) Delete(ctx context.Context, params db.DeleteSuscripcionParams) error {
	return r.queries.DeleteSuscripcion(ctx, params)
}
