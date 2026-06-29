package repository

import (
    "database/sql"
    "fmt"
    "time"
    
    "github.com/lib/pq"
    "bank-service/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    
    now := time.Now()
    err := r.db.QueryRow(
        query,
        user.Username,
        user.Email,
        user.PasswordHash,
        now,
        now,
    ).Scan(&user.ID)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" {
                return fmt.Errorf("user with this email or username already exists")
            }
        }
        return err
    }
    
    user.CreatedAt = now
    user.UpdatedAt = now
    return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, username, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    
    var user models.User
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}
