package service

import (
	"context"
	"errors"
	"go_finance/internal/api/middleware"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"go_finance/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shopspring/decimal"
)

type userService struct {
	userRepo        repository.UserRepository
	balanceRepo     repository.BalanceRepository
	jwtSecret       string
	jwtExpiresIn    time.Duration
	auditLogService AuditLogService
}

func NewUserService(
	repo repository.UserRepository,
	balanceRepo repository.BalanceRepository,
	secret string,
	expiresIn time.Duration,
	auditLogService AuditLogService,
) *userService {
	return &userService{
		userRepo:        repo,
		balanceRepo:     balanceRepo,
		jwtSecret:       secret,
		jwtExpiresIn:    expiresIn,
		auditLogService: auditLogService,
	}
}

func (s *userService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	auditDetails := "User registration attempt: username=" + req.Username + ", email=" + req.Email
	s.auditLogService.Create(ctx, "user", "", "register_attempt", auditDetails)
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	newUser := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := newUser.HashPassword(req.Password); err != nil {
		utils.Logger.Error("Error hashing password", "error", err)
		s.auditLogService.Create(ctx, "user", "", "register_failed", "Password hashing failed")
		return nil, errors.New("could not process request")
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		utils.Logger.Error("Error creating user in repository", "error", err)
		s.auditLogService.Create(ctx, "user", "", "register_failed", "User creation failed")
		return nil, errors.New("could not create user")
	}

	amount, _ := decimal.NewFromString("0.00")
	if err := s.balanceRepo.CreateBalance(ctx, &domain.Balance{
		UserID:        newUser.ID,
		Amount:        amount,
		LastUpdatedAt: time.Now(),
	}); err != nil {
		utils.Logger.Error("Error creating user's balance", "error", err)
		s.auditLogService.Create(ctx, "balance", newUser.ID, "balance_create_failed", "Balance creation failed for user")
	}
	s.auditLogService.Create(ctx, "user", newUser.ID, "register_success", "User registered successfully")

	return &RegisterResponse{
		User:    newUser,
		Message: "User registered successfully",
	}, nil
}

func (s *userService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	s.auditLogService.Create(ctx, "user", "", "login_attempt", "Login attempt for email="+req.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		utils.Logger.Warn("User not found by email", "email", req.Email, "error", err)
		s.auditLogService.Create(ctx, "user", "", "login_failed", "User not found for email="+req.Email)
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		s.auditLogService.Create(ctx, "user", user.ID, "login_failed", "Invalid password for user")
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		utils.Logger.Error("Could not generate JWT", "error", err)
		s.auditLogService.Create(ctx, "user", user.ID, "login_failed", "JWT generation failed")
		return nil, errors.New("could not process login")
	}
	s.auditLogService.Create(ctx, "user", user.ID, "login_success", "User logged in successfully")

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
	s.auditLogService.Create(ctx, "user", req.ID, "update_attempt", "User update attempt")
	if ctx.Value(middleware.UserRoleKey) != "admin" {
		if ctx.Value(middleware.UserIDKey) != req.ID {
			return nil, errors.New("you dont have any access")
		}
	}

	user, err := s.userRepo.UpdateUserByID(ctx, req.ID)
	if err != nil {
		s.auditLogService.Create(ctx, "user", req.ID, "update_failed", "User update failed")
		return nil, err
	}
	s.auditLogService.Create(ctx, "user", req.ID, "update_success", "User updated successfully")

	return &PutUserByIdResponse{
		User:    user,
		Message: "User Updated Succesfully",
	}, nil
}
func (s *userService) DeleteUser(ctx context.Context, req DeleteUserByIdRequest) (*DeleteUserByIdResponse, error) {
	s.auditLogService.Create(ctx, "user", req.ID, "delete_attempt", "User delete attempt")
	if ctx.Value(middleware.UserRoleKey) != "admin" {
		if ctx.Value(middleware.UserIDKey) != req.ID {
			return nil, errors.New("you dont have any access")
		}
	}

	_, err := s.userRepo.DeleteUserByID(ctx, req.ID)
	if err != nil {
		s.auditLogService.Create(ctx, "user", req.ID, "delete_failed", "User delete failed")
		return nil, err
	}
	s.auditLogService.Create(ctx, "user", req.ID, "delete_success", "User deleted successfully")

	return &DeleteUserByIdResponse{
		Message: "User deleted successfully",
	}, nil
}
