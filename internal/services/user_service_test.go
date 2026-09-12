package services

import (
	"context"
	"testing"
)

// MockUserRepository simula el comportamiento del repositorio sin tocar Postgres
type MockUserRepository struct{}

func (m *MockUserRepository) Create(ctx context.Context, email string) error {
	return nil
}

func TestUserService_RegisterUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Email válido", "test@example.com", false},
		{"Email inválido", "email-malformado", true},
		{"Email vacío", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RegisterUser(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}