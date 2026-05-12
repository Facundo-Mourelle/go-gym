package service

import (
	"fmt"
	"log"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/config"
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo        repository.UserRepository
	routineService  *RoutineService
	jwtSecret       []byte
	jwtExpiry       time.Duration
}

func NewAuthService(userRepo repository.UserRepository, routineService *RoutineService, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		routineService:  routineService,
		jwtSecret:       []byte(cfg.JWTSecret),
		jwtExpiry:       time.Duration(cfg.JWTExpirationHours) * time.Hour,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req RegisterRequest) (AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return AuthResponse{}, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := domain.User{
		ID:           domain.UserID(generateID()),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return AuthResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.routineService.SeedStarterRoutines(string(user.ID)); err != nil {
		log.Printf("Failed to seed starter routines for user %s: %v", user.ID, err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{
		Token:  token,
		UserID: string(user.ID),
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(*user)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return AuthResponse{
		Token:  token,
		UserID: string(user.ID),
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return "", fmt.Errorf("invalid token expiration")
		}

		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", fmt.Errorf("token expired")
		}

		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}

func (s *AuthService) generateToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": string(user.ID),
		"email":   user.Email,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Request/Response Types

type RegisterRequest struct {
	Email    string
	Password string
	Name     string
}

type LoginRequest struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
