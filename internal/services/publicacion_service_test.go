package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "emprendimientos.com/servidor-go/db/sqlc"
)

// mockPublicacionRepo simula el repositorio de publicaciones para pruebas unitarias aisladas
type mockPublicacionRepo struct {
	createFunc func(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error)
	getByIDFunc func(ctx context.Context, id int32) (db.Publicacione, error)
	listByEmprendimientoFunc func(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error)
	updateFunc func(ctx context.Context, params db.UpdatePublicacionParams) error
	deleteFunc func(ctx context.Context, id int32) error
}

func (m *mockPublicacionRepo) Create(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, params)
	}
	return db.Publicacione{}, nil
}

func (m *mockPublicacionRepo) GetByID(ctx context.Context, id int32) (db.Publicacione, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return db.Publicacione{IDPublicacion: id}, nil
}

func (m *mockPublicacionRepo) ListByEmprendimiento(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error) {
	if m.listByEmprendimientoFunc != nil {
		return m.listByEmprendimientoFunc(ctx, idEmprendimiento)
	}
	return []db.Publicacione{}, nil
}

func (m *mockPublicacionRepo) Update(ctx context.Context, params db.UpdatePublicacionParams) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, params)
	}
	return nil
}

func (m *mockPublicacionRepo) Delete(ctx context.Context, id int32) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func floatPtr(v float64) *float64 {
	return &v
}

func TestPublicacionService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("falla con ID de emprendimiento invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		input := CreatePublicacionInput{
			IDEmprendimiento: 0,
			Titulo:           "Oferta valida",
			Tipo:             "promocion",
		}

		_, err := service.Create(ctx, input)
		if !errors.Is(err, ErrIDEmprendimientoInvalido) {
			t.Errorf("se esperaba ErrIDEmprendimientoInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("falla con titulo vacio o de puros espacios", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		input := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "   ",
			Tipo:             "producto",
		}

		_, err := service.Create(ctx, input)
		if !errors.Is(err, ErrTituloRequerido) {
			t.Errorf("se esperaba ErrTituloRequerido, se obtuvo: %v", err)
		}
	})

	t.Run("falla con tipo invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		input := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "Titulo correcto",
			Tipo:             "invalido",
		}

		_, err := service.Create(ctx, input)
		if !errors.Is(err, ErrTipoInvalido) {
			t.Errorf("se esperaba ErrTipoInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("falla con precio negativo o cero", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})

		// Precio 0
		inputCero := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "Titulo correcto",
			Tipo:             "servicio",
			Precio:           floatPtr(0),
		}
		_, err := service.Create(ctx, inputCero)
		if !errors.Is(err, ErrPrecioInvalido) {
			t.Errorf("con precio 0 se esperaba ErrPrecioInvalido, se obtuvo: %v", err)
		}

		// Precio negativo
		inputNegativo := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "Titulo correcto",
			Tipo:             "servicio",
			Precio:           floatPtr(-50.0),
		}
		_, err = service.Create(ctx, inputNegativo)
		if !errors.Is(err, ErrPrecioInvalido) {
			t.Errorf("con precio negativo se esperaba ErrPrecioInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("exito con precio nulo", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			createFunc: func(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error) {
				if params.Precio.Valid {
					t.Errorf("se esperaba precio nulo (Valid: false), se obtuvo: %+v", params.Precio)
				}
				if params.Titulo != "Titulo Valido" {
					t.Errorf("titulo esperado 'Titulo Valido', obtenido '%s'", params.Titulo)
				}
				return db.Publicacione{
					IDPublicacion:    10,
					IDEmprendimiento: params.IDEmprendimiento,
					Titulo:           params.Titulo,
					Tipo:             params.Tipo,
				}, nil
			},
		}

		service := NewPublicacionService(mockRepo)
		input := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "  Titulo Valido  ",
			Tipo:             "OTRO", // Debe normalizarse a minusculas
			Precio:           nil,
		}

		pub, err := service.Create(ctx, input)
		if err != nil {
			t.Fatalf("no se esperaba error, se obtuvo: %v", err)
		}
		if pub.IDPublicacion != 10 {
			t.Errorf("ID esperado 10, obtenido %d", pub.IDPublicacion)
		}
	})

	t.Run("exito con precio positivo y formateado", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			createFunc: func(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error) {
				if !params.Precio.Valid || params.Precio.String != "1250.50" {
					t.Errorf("precio esperado '1250.50', obtenido '%s'", params.Precio.String)
				}
				return db.Publicacione{
					IDPublicacion: 20,
					Titulo:        params.Titulo,
				}, nil
			},
		}

		service := NewPublicacionService(mockRepo)
		input := CreatePublicacionInput{
			IDEmprendimiento: 2,
			Titulo:           "Producto Especial",
			Tipo:             "producto",
			Precio:           floatPtr(1250.5),
		}

		pub, err := service.Create(ctx, input)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if pub.IDPublicacion != 20 {
			t.Errorf("esperado 20, obtenido %d", pub.IDPublicacion)
		}
	})

	t.Run("propaga error de repositorio", func(t *testing.T) {
		expectedErr := errors.New("db connection failure")
		mockRepo := &mockPublicacionRepo{
			createFunc: func(ctx context.Context, params db.CreatePublicacionParams) (db.Publicacione, error) {
				return db.Publicacione{}, expectedErr
			},
		}

		service := NewPublicacionService(mockRepo)
		input := CreatePublicacionInput{
			IDEmprendimiento: 1,
			Titulo:           "Producto",
			Tipo:             "producto",
		}

		_, err := service.Create(ctx, input)
		if err == nil {
			t.Fatal("se esperaba error de repositorio pero fue nil")
		}
	})
}

func TestPublicacionService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("falla con ID menor o igual a cero", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})

		_, err := service.GetByID(ctx, 0)
		if !errors.Is(err, ErrIDPublicacionInvalido) {
			t.Errorf("se esperaba ErrIDPublicacionInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("retorna ErrPublicacionNoEncontrada cuando no existe", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{}, sql.ErrNoRows
			},
		}

		service := NewPublicacionService(mockRepo)
		_, err := service.GetByID(ctx, 999)
		if !errors.Is(err, ErrPublicacionNoEncontrada) {
			t.Errorf("se esperaba ErrPublicacionNoEncontrada, se obtuvo: %v", err)
		}
	})

	t.Run("exito al obtener publicacion", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{
					IDPublicacion: id,
					Titulo:        "Publicacion Encontrada",
				}, nil
			},
		}

		service := NewPublicacionService(mockRepo)
		pub, err := service.GetByID(ctx, 5)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if pub.Titulo != "Publicacion Encontrada" {
			t.Errorf("titulo inesperado: %s", pub.Titulo)
		}
	})
}

func TestPublicacionService_ListByEmprendimiento(t *testing.T) {
	ctx := context.Background()

	t.Run("falla con ID de emprendimiento invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		_, err := service.ListByEmprendimiento(ctx, -1)
		if !errors.Is(err, ErrIDEmprendimientoInvalido) {
			t.Errorf("se esperaba ErrIDEmprendimientoInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("exito al listar publicaciones", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			listByEmprendimientoFunc: func(ctx context.Context, idEmprendimiento int32) ([]db.Publicacione, error) {
				return []db.Publicacione{
					{IDPublicacion: 1, Titulo: "Pub 1"},
					{IDPublicacion: 2, Titulo: "Pub 2"},
				}, nil
			},
		}

		service := NewPublicacionService(mockRepo)
		pubs, err := service.ListByEmprendimiento(ctx, 10)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if len(pubs) != 2 {
			t.Errorf("se esperaban 2 publicaciones, se obtuvieron %d", len(pubs))
		}
	})
}

func TestPublicacionService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("falla con ID invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		err := service.Update(ctx, 0, UpdatePublicacionInput{Titulo: "Test", Tipo: "otro"})
		if !errors.Is(err, ErrIDPublicacionInvalido) {
			t.Errorf("se esperaba ErrIDPublicacionInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("falla cuando publicacion no existe", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{}, sql.ErrNoRows
			},
		}

		service := NewPublicacionService(mockRepo)
		err := service.Update(ctx, 50, UpdatePublicacionInput{Titulo: "Test", Tipo: "otro"})
		if !errors.Is(err, ErrPublicacionNoEncontrada) {
			t.Errorf("se esperaba ErrPublicacionNoEncontrada, se obtuvo: %v", err)
		}
	})

	t.Run("falla con titulo vacio", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		err := service.Update(ctx, 1, UpdatePublicacionInput{Titulo: "  ", Tipo: "otro"})
		if !errors.Is(err, ErrTituloRequerido) {
			t.Errorf("se esperaba ErrTituloRequerido, se obtuvo: %v", err)
		}
	})

	t.Run("falla con tipo invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		err := service.Update(ctx, 1, UpdatePublicacionInput{Titulo: "Valido", Tipo: "categoria_inexistente"})
		if !errors.Is(err, ErrTipoInvalido) {
			t.Errorf("se esperaba ErrTipoInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("falla con precio invalido", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		err := service.Update(ctx, 1, UpdatePublicacionInput{
			Titulo: "Valido",
			Tipo:   "servicio",
			Precio: floatPtr(-5),
		})
		if !errors.Is(err, ErrPrecioInvalido) {
			t.Errorf("se esperaba ErrPrecioInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("exito en update", func(t *testing.T) {
		updatedCalled := false
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{IDPublicacion: id}, nil
			},
			updateFunc: func(ctx context.Context, params db.UpdatePublicacionParams) error {
				updatedCalled = true
				if params.IDPublicacion != 10 {
					t.Errorf("esperado ID 10, obtenido %d", params.IDPublicacion)
				}
				if params.Titulo != "Titulo Nuevo" {
					t.Errorf("esperado 'Titulo Nuevo', obtenido '%s'", params.Titulo)
				}
				return nil
			},
		}

		service := NewPublicacionService(mockRepo)
		err := service.Update(ctx, 10, UpdatePublicacionInput{
			Titulo: "Titulo Nuevo",
			Tipo:   "producto",
			Precio: floatPtr(500),
		})
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if !updatedCalled {
			t.Errorf("se esperaba que repo.Update fuera llamado")
		}
	})
}

func TestPublicacionService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("falla con ID menor o igual a cero", func(t *testing.T) {
		service := NewPublicacionService(&mockPublicacionRepo{})
		err := service.Delete(ctx, 0)
		if !errors.Is(err, ErrIDPublicacionInvalido) {
			t.Errorf("se esperaba ErrIDPublicacionInvalido, se obtuvo: %v", err)
		}
	})

	t.Run("falla cuando publicacion no existe", func(t *testing.T) {
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{}, sql.ErrNoRows
			},
		}

		service := NewPublicacionService(mockRepo)
		err := service.Delete(ctx, 123)
		if !errors.Is(err, ErrPublicacionNoEncontrada) {
			t.Errorf("se esperaba ErrPublicacionNoEncontrada, se obtuvo: %v", err)
		}
	})

	t.Run("exito al eliminar", func(t *testing.T) {
		deleteCalled := false
		mockRepo := &mockPublicacionRepo{
			getByIDFunc: func(ctx context.Context, id int32) (db.Publicacione, error) {
				return db.Publicacione{IDPublicacion: id}, nil
			},
			deleteFunc: func(ctx context.Context, id int32) error {
				deleteCalled = true
				return nil
			},
		}

		service := NewPublicacionService(mockRepo)
		err := service.Delete(ctx, 123)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if !deleteCalled {
			t.Errorf("se esperaba que repo.Delete fuera llamado")
		}
	})
}
