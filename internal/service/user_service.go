package service

import (
	"context"
	"errors"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userService struct {
	userRepo     repository.UserRepository
	jwtSecret    string
	jwtExpiresIn time.Duration
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository, secret string, expiresIn time.Duration) *userService {
	return &userService{
		userRepo:     repo,
		jwtSecret:    secret,
		jwtExpiresIn: expiresIn,
	}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// E-posta adresinin daha önce alınıp alınmadığını kontrol et
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	newUser := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := newUser.HashPassword(req.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("could not process request")
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		log.Printf("Error creating user in repo: %v", err)
		return nil, errors.New("could not create user")
	}

	return &RegisterResponse{
		User:    newUser,
		Message: "User registered successfully",
	}, nil
}

func (s *userService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("User not found by email %s: %v", req.Email, err)
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// JWT Oluştur
	token, err := s.generateJWT(user)
	if err != nil {
		log.Printf("Could not generate JWT: %v", err)
		return nil, errors.New("could not process login")
	}

	return &LoginResponse{
		User:        user,
		AccessToken: token,
		Message:     "Login successful",
	}, nil
}

func (s *userService) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.jwtExpiresIn).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
