package models

import "time"

type Credit struct {
    ID              int       `json:"id" db:"id"`
    UserID          int       `json:"user_id" db:"user_id"`
    AccountID       int       `json:"account_id" db:"account_id"`
    Amount          float64   `json:"amount" db:"amount"`
    RemainingAmount float64   `json:"remaining_amount" db:"remaining_amount"`
    InterestRate    float64   `json:"interest_rate" db:"interest_rate"`
    TermMonths      int       `json:"term_months" db:"term_months"`
    MonthlyPayment  float64   `json:"monthly_payment" db:"monthly_payment"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreditRequest struct {
    Amount     float64 `json:"amount" validate:"required,gt=0"`
    TermMonths int     `json:"term_months" validate:"required,min=3,max=60"`
    AccountID  int     `json:"account_id" validate:"required"`
}

type PaymentSchedule struct {
    ID          int       `json:"id" db:"id"`
    CreditID    int       `json:"credit_id" db:"credit_id"`
    PaymentDate time.Time `json:"payment_date" db:"payment_date"`
    Amount      float64   `json:"amount" db:"amount"`
    Principal   float64   `json:"principal" db:"principal"`
    Interest    float64   `json:"interest" db:"interest"`
    Status      string    `json:"status" db:"status"`
    Penalty     float64   `json:"penalty" db:"penalty"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
