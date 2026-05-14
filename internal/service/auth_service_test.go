package service

import (
	"testing"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/config"
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository/memory"
)

func TestRegister(t *testing.T) {
	userRepo := memory.NewUserMemoryRepository()
	routineRepo := memory.NewRoutineMemoryRepository()
	routineService := NewRoutineService(routineRepo)
	
	cfg := &config.Config{
		JWTSecret:          "test-secret-key",
		JWTExpirationHours: 24,
	}
	
	authService := NewAuthService(userRepo, routineService, cfg)

	t.Run("successful registration", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "SecurePass123!",
			Name:     "Test User",
		}

		resp, err := authService.Register(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Email != req.Email {
			t.Errorf("expected email %s, got %s", req.Email, resp.Email)
		}

		if resp.Token == "" {
			t.Error("expected token to be generated")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "duplicate@example.com",
			Password: "SecurePass123!",
			Name:     "User One",
		}

		_, err := authService.Register(req)
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		_, err = authService.Register(req)
		if err == nil {
			t.Error("expected error for duplicate email, got nil")
		}
	})
}

func TestLogin(t *testing.T) {
	userRepo := memory.NewUserMemoryRepository()
	routineRepo := memory.NewRoutineMemoryRepository()
	routineService := NewRoutineService(routineRepo)
	
	cfg := &config.Config{
		JWTSecret:          "test-secret-key",
		JWTExpirationHours: 24,
	}
	
	authService := NewAuthService(userRepo, routineService, cfg)

	email := "login@example.com"
	password := "SecurePass123!"

	_, err := authService.Register(RegisterRequest{
		Email:    email,
		Password: password,
		Name:     "Login User",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	t.Run("successful login", func(t *testing.T) {
		resp, err := authService.Login(LoginRequest{
			Email:    email,
			Password: password,
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Email != email {
			t.Errorf("expected email %s, got %s", email, resp.Email)
		}

		if resp.Token == "" {
			t.Error("expected token to be generated")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		_, err := authService.Login(LoginRequest{
			Email:    email,
			Password: "WrongPassword123!",
		})

		if err == nil {
			t.Error("expected error for invalid password, got nil")
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := authService.Login(LoginRequest{
			Email:    "nonexistent@example.com",
			Password: password,
		})

		if err == nil {
			t.Error("expected error for non-existent user, got nil")
		}
	})
}

func TestValidateToken(t *testing.T) {
	userRepo := memory.NewUserMemoryRepository()
	routineRepo := memory.NewRoutineMemoryRepository()
	routineService := NewRoutineService(routineRepo)
	
	cfg := &config.Config{
		JWTSecret:          "test-secret-key",
		JWTExpirationHours: 24,
	}
	
	authService := NewAuthService(userRepo, routineService, cfg)

	resp, err := authService.Register(RegisterRequest{
		Email:    "token@example.com",
		Password: "SecurePass123!",
		Name:     "Token User",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	t.Run("valid token", func(t *testing.T) {
		userID, err := authService.ValidateToken(resp.Token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if userID != resp.UserID {
			t.Errorf("expected userID %s, got %s", resp.UserID, userID)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := authService.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expiredCfg := &config.Config{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: -1,
		}
		expiredAuthService := NewAuthService(userRepo, routineService, expiredCfg)

		user := domain.User{
			ID:           domain.UserID("test-user"),
			Email:        "expired@example.com",
			PasswordHash: "hash",
			Name:         "Expired User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		token, err := expiredAuthService.generateToken(user)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		_, err = authService.ValidateToken(token)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})
}
