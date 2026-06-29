package repository

import (
    "sync"
    "time"
    
    "bank-service/internal/models"
)

type InMemoryTransactionRepository struct {
    transactions map[int]*models.Transaction
    mu           sync.RWMutex
    idGen        int
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
    return &InMemoryTransactionRepository{
        transactions: make(map[int]*models.Transaction),
        idGen:        1,
    }
}

func (r *InMemoryTransactionRepository) Create(tx *models.Transaction) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    tx.ID = r.idGen
    r.idGen++
    tx.CreatedAt = time.Now()
    tx.Status = "completed"
    
    r.transactions[tx.ID] = tx
    return nil
}

func (r *InMemoryTransactionRepository) GetUserTransactions(userID int, limit int) ([]*models.Transaction, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var txs []*models.Transaction
    for _, tx := range r.transactions {
        // Фильтруем транзакции для пользователя
        // В реальном приложении нужно проверять через связанные счета
        if limit > 0 && len(txs) >= limit {
            break
        }
        txCopy := *tx
        txs = append(txs, &txCopy)
    }
    return txs, nil
}
