package models

import "time"

type Card struct {
    ID               int       `json:"id" db:"id"`
    UserID           int       `json:"user_id" db:"user_id"`
    AccountID        int       `json:"account_id" db:"account_id"`
    NumberEncrypted  string    `json:"-" db:"number_encrypted"`
    NumberHMAC       string    `json:"-" db:"number_hmac"`
    ExpiryEncrypted  string    `json:"-" db:"expiry_encrypted"`
    CVVHash          string    `json:"-" db:"cvv_hash"`
    CardType         string    `json:"card_type" db:"card_type"`
    Status           string    `json:"status" db:"status"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type CardResponse struct {
    ID       int    `json:"id"`
    Number   string `json:"number"`
    Expiry   string `json:"expiry"`
    CardType string `json:"card_type"`
    Status   string `json:"status"`
}

type CreateCardRequest struct {
    AccountID int `json:"account_id" validate:"required"`
}

type PaymentRequest struct {
    CardNumber string  `json:"card_number" validate:"required"`
    CVV        string  `json:"cvv" validate:"required,len=3"`
    Amount     float64 `json:"amount" validate:"required,gt=0"`
}
