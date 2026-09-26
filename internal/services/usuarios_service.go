package services

//logica de negocio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	db "emprendimientos.com/servidor-go/db/sqlc"
	"emprendimientos.com/servidor-go/internal/repositories"

)

// variables de errores para moestrar segun el error
var (
	ErrFormatoInvalidoEmail = errors.New("formato de EMAIL invalido.")
	ErrNombreVacio      = errors.New("el nombre completo es OBLIGATORIO")
	ErrClave           = errors.New("la CLAVE debe tener al menos 6 caracteres")
	ErrRolInvalido          = errors.New("el ROL debe ser 'cliente', 'emprendedor' o 'administrador'")
	ErrEmailDuplicado       = errors.New("el EMAIL ya se encuentra registrado")
	ErrUsuarioNoEncontrado  = errors.New("usuario NO encontrado")
)

//FUNCIONES ADICIONALES



//validar formatos -> email (teoria)
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email){
		return  errors.New("FORMATO DE EMAIL INVALIDO.")
	}
	return nil
}

func validarUsuario(nombre, email, rol string) error{
	
	//REGLA 1  
	if strings.TrimSpace(nombre) == "" {
		return ErrNombreVacio
	}

	//REGLA 2 
	if err:= ValidateEmail(strings.TrimSpace(email)); err != nil {
		return fmt.Errorf("Error al validar el email: %w", err)
	}



	rolNormalizado := strings.TrimSpace(strings.ToLower(rol))
	if rolNormalizado != "cliente" && rolNormalizado != "emprendedor" && rolNormalizado != "administrador" {
		return ErrRolInvalido
	}

	return nil

}

// UserService  es la logica de negocio para que usen los usuarios
type UserService interface {
	CrearUsuario(ctx context.Context, nombre, email, clave, rol string, idEmprendimiento *int32) (db.Usuario, error)
	ObtenerUsuario(ctx context.Context, id int32) (db.Usuario, error)
	ObtenerUsuarioPorEmail(ctx context.Context, email string) (db.Usuario, error)
	ListarUsuarios(ctx context.Context) ([]db.Usuario, error)
	ActualizarUsuario(ctx context.Context, id int32, nombre, email, clave, rol string, idEmprendimiento *int32) error
	EliminarUsuario(ctx context.Context, id int32) error
}

//convencion de catedra, el struct en minuscula
type userService struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) UserService {
	return &userService{repo: repo}
}



// BUSCAR POR EMAIL

// ObtenerUsuarioPorEmail implementa la búsqueda con validación previa de formato
func (s *userService) ObtenerUsuarioPorEmail(ctx context.Context, email string) (db.Usuario, error) {
	emailLimpio := strings.TrimSpace(strings.ToLower(email))
	if err := ValidateEmail(emailLimpio); err != nil {
		return db.Usuario{}, fmt.Errorf("email inválido: %w", err)
	}

	usuario, err := s.repo.GetByEmail(ctx, emailLimpio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Usuario{}, ErrUsuarioNoEncontrado
		}
		return db.Usuario{}, fmt.Errorf("error al buscar usuario por email: %w", err)
	}

	return usuario, nil
}

//crearu usuario

func (s* userService) CrearUsuario(ctx context.Context, nombre, email, clave, rol string, idEmprendimiento *int32) (db.Usuario, error){

	if err := validarUsuario(nombre, email, rol); err != nil {
		return db.Usuario{}, err
	}

	if len(clave) < 6 {
		return db.Usuario{}, ErrClave
	}

	//verificar email unico
	emailNormalizado := strings.TrimSpace(strings.ToLower(email))

	
	_, err := s.repo.GetByEmail(ctx, emailNormalizado)
		if err == nil {
			return db.Usuario{}, ErrEmailDuplicado
		} else if !errors.Is(err, sql.ErrNoRows) {
			return db.Usuario{}, fmt.Errorf("error al verificar unicidad de email: %w", err)
		}

	params := db.CreateUsuarioParams{
		NombreCompleto: strings.TrimSpace(nombre),
		Email:          strings.TrimSpace(strings.ToLower(email)),
		Clave:          clave,
		Rol:            strings.TrimSpace(strings.ToLower(rol)),
		IDEmprendimiento: sql.NullInt32{
			Valid: idEmprendimiento != nil,
		},
	}

	if idEmprendimiento != nil {
		params.IDEmprendimiento.Int32 = *idEmprendimiento
	}

	// manejo de persistencia mediante el repo
	return s.repo.Create(ctx, params)

}

// Obtener usuario

func (s *userService) ObtenerUsuario(ctx context.Context, id int32) (db.Usuario, error){
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListarUsuarios(ctx context.Context) ([]db.Usuario, error){
	return s.repo.ListAll(ctx)
}

func (s *userService) ActualizarUsuario(ctx context.Context, id int32, nombre, email, clave, rol string, idEmprendimiento *int32) error{
	if err := validarUsuario(nombre, email, rol); err != nil {
		return err
	}

	// Si en el update también se exige clave:
	if len(clave) < 6 {
		return ErrClave
	}

	params := db.UpdateUsuarioParams{
		IDUsuario:      id,
		NombreCompleto: strings.TrimSpace(nombre),
		Email:          strings.TrimSpace(strings.ToLower(email)),
		Clave:          clave,
		Rol:            strings.TrimSpace(strings.ToLower(rol)),
		IDEmprendimiento: sql.NullInt32{
			Valid: idEmprendimiento != nil,
		},
	}
	if idEmprendimiento != nil {
		params.IDEmprendimiento.Int32 = *idEmprendimiento
	}

	return s.repo.Update(ctx, params)

}

func (s *userService) EliminarUsuario(ctx context.Context, id int32) error {
		return s.repo.Delete(ctx, id)
}