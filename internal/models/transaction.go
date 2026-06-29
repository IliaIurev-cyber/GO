package models

import "time"

type Transaction struct {
    ID            int       `json:"id" db:"id"`
    FromAccountID *int      `json:"from_account_id" db:"from_account_id"`
    ToAccountID   *int      `json:"to_account_id" db:"to_account_id"`
    Amount        float64   `json:"amount" db:"amount"`
    Type          string    `json:"type" db:"type"`
    Status        string    `json:"status" db:"status"`
    Description   string    `json:"description" db:"description"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type TransferRequest struct {
    FromAccountID int     `json:"from_account_id" validate:"required"`
    ToAccountID   int     `json:"to_account_id" validate:"required"`
    Amount        float64 `json:"amount" validate:"required,gt=0"`
    Description   string  `json:"description"`
}

type DepositRequest struct {
    AccountID int     `json:"account_id" validate:"required"`
    Amount    float64 `json:"amount" validate:"required,gt=0"`
}
