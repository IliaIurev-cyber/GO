package models

import "time"

type Account struct {
    ID        int       `json:"id" db:"id"`
    UserID    int       `json:"user_id" db:"user_id"`
    Number    string    `json:"number" db:"number"`
    Balance   float64   `json:"balance" db:"balance"`
    Currency  string    `json:"currency" db:"currency"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
    Currency string `json:"currency" validate:"required,oneof=RUB"`
}

type AccountResponse struct {
    ID        int     `json:"id"`
    Number    string  `json:"number"`
    Balance   float64 `json:"balance"`
    Currency  string  `json:"currency"`
    Status    string  `json:"status"`
}
