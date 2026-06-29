package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "bank-service/internal/service"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }
            
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == authHeader {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }
            
            userID, err := authService.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserID(ctx context.Context) int {
    if userID, ok := ctx.Value(UserIDKey).(int); ok {
        return userID
    }
    return 0
}
