package repository

import (
    "fmt"
    "sync"
    "time"
    
    "bank-service/internal/models"
)

type InMemoryUserRepository struct {
    users map[string]*models.User // key: email
    mu    sync.RWMutex
    idGen int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*models.User),
        idGen: 1,
    }
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Check if user exists
    if _, exists := r.users[user.Email]; exists {
        return fmt.Errorf("user with this email already exists")
    }
    
    // Set ID and timestamps
    user.ID = r.idGen
    r.idGen++
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now
    
    // Store user
    r.users[user.Email] = user
    
    return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, exists := r.users[email]
    if !exists {
        return nil, nil
    }
    
    // Return a copy to avoid modification
    userCopy := *user
    return &userCopy, nil
}
