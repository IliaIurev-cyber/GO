package repository

import (
    "fmt"
    "sync"
    "time"
    
    "bank-service/internal/models"
)

type InMemoryAccountRepository struct {
    accounts map[int]*models.Account
    mu       sync.RWMutex
    idGen    int
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
    return &InMemoryAccountRepository{
        accounts: make(map[int]*models.Account),
        idGen:    1,
    }
}

func (r *InMemoryAccountRepository) Create(account *models.Account) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    account.ID = r.idGen
    r.idGen++
    now := time.Now()
    account.CreatedAt = now
    account.UpdatedAt = now
    account.Status = "active"
    
    r.accounts[account.ID] = account
    return nil
}

func (r *InMemoryAccountRepository) FindByID(id int) (*models.Account, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    account, exists := r.accounts[id]
    if !exists {
        return nil, nil
    }
    
    accountCopy := *account
    return &accountCopy, nil
}

func (r *InMemoryAccountRepository) FindByUserID(userID int) ([]*models.Account, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var accounts []*models.Account
    for _, account := range r.accounts {
        if account.UserID == userID && account.Status == "active" {
            accountCopy := *account
            accounts = append(accounts, &accountCopy)
        }
    }
    return accounts, nil
}

func (r *InMemoryAccountRepository) UpdateBalance(id int, amount float64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    account, exists := r.accounts[id]
    if !exists {
        return fmt.Errorf("account not found")
    }
    if account.Status != "active" {
        return fmt.Errorf("account is not active")
    }
    
    account.Balance += amount
    account.UpdatedAt = time.Now()
    return nil
}
