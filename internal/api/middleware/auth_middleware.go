package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

type contextKey int

const (
    UserIDKey contextKey = iota
    UserRoleKey 
    UserEmailKey
    UserExpKey  
)

func (k contextKey) String() string {
    switch k {
    case UserIDKey:
        return "user_id"
    case UserRoleKey:
        return "user_role"
    case UserEmailKey:
        return "user_email"
    case UserExpKey:
        return "user_exp"
    default:
        return "unknown"
    }
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: secret,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := bearerToken[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			http.Error(w, "Token expired", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims["sub"])
		ctx = context.WithValue(ctx, UserRoleKey, claims["role"])
		ctx = context.WithValue(ctx, UserEmailKey, claims["email"])
		ctx = context.WithValue(ctx, UserExpKey, claims["exp"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
