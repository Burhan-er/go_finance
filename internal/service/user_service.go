package service

import (
	"context"
	"errors"
	"go_finance/internal/api/middleware"
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

// GetUserByID retrieves a single user by their ID.
func (s *userService) GetUserByID(ctx context.Context, req *GetUserByIdRequest) (*GetUserByIdResponse, error) {

	if ctx.Value(middleware.UserRoleKey) != "admin" {
		if ctx.Value(middleware.UserIDKey) != req.ID {
			return nil, errors.New("you dont have any access")
		}
	}

	user, err := s.userRepo.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserByIdResponse{
		User: user,
	}, nil
}

func (s *userService) GetAllUsers(ctx context.Context, req GetAllUsersRequest) (*GetAllUsersResponse, error) {
	if ctx.Value(middleware.UserRoleKey) != "admin" {
		return nil, errors.New("you dont have any access")
	}
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &GetAllUsersResponse{
		Users: users,
	}, nil
}
func (s *userService) UpdateUser(ctx context.Context, req PutUserByIdRequest) (*PutUserByIdResponse, error) {
	if ctx.Value(middleware.UserRoleKey) != "admin" {
		if ctx.Value(middleware.UserIDKey) != req.ID {
			return nil, errors.New("you dont have any access")
		}
	}

	user, err := s.userRepo.UpdateUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &PutUserByIdResponse{
		User:    user,
		Message: "User Updated Succesfully",
	}, nil
}
func (s *userService) DeleteUser(ctx context.Context, req DeleteUserByIdRequest) (*DeleteUserByIdResponse, error) {
	if ctx.Value(middleware.UserRoleKey) != "admin" {
		if ctx.Value(middleware.UserIDKey) != req.ID {
			return nil, errors.New("you dont have any access")
		}
	}

	_, err := s.userRepo.DeleteUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteUserByIdResponse{
		Message: "User deleted successfully",
	}, nil
}
