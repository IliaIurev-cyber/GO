package service

import (
    "errors"
    "strconv"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "bank-service/internal/config"
    "bank-service/internal/models"
    "bank-service/internal/utils"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
}

type AuthService struct {
    userRepo UserRepository
    config   *config.Config
}

func NewAuthService(userRepo UserRepository, config *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        config:   config,
    }
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
    // Check if user exists
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, *models.User, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return "", nil, err
    }
    if user == nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
        return "", nil, errors.New("invalid credentials")
    }
    
    token, err := s.generateJWT(user)
    if err != nil {
        return "", nil, err
    }
    
    return token, user, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   strconv.Itoa(user.ID),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.config.JWTSecret), nil
    })
    
    if err != nil || !token.Valid {
        return 0, errors.New("invalid token")
    }
    
    userID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        return 0, errors.New("invalid user ID in token")
    }
    
    return userID, nil
}
