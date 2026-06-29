package service

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
    
    "bank-service/internal/models"
    "bank-service/internal/repository"
)

type AccountService struct {
    accountRepo     *repository.InMemoryAccountRepository
    transactionRepo *repository.InMemoryTransactionRepository
}

func NewAccountService(
    accountRepo *repository.InMemoryAccountRepository,
    transactionRepo *repository.InMemoryTransactionRepository,
) *AccountService {
    return &AccountService{
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
    }
}

func (s *AccountService) CreateAccount(userID int, req *models.CreateAccountRequest) (*models.Account, error) {
    account := &models.Account{
        UserID:   userID,
        Number:   s.generateAccountNumber(),
        Balance:  0,
        Currency: req.Currency,
    }
    
    if err := s.accountRepo.Create(account); err != nil {
        return nil, err
    }
    
    return account, nil
}

func (s *AccountService) GetAccounts(userID int) ([]*models.Account, error) {
    return s.accountRepo.FindByUserID(userID)
}

func (s *AccountService) GetAccount(userID, accountID int) (*models.Account, error) {
    account, err := s.accountRepo.FindByID(accountID)
    if err != nil {
        return nil, err
    }
    if account == nil {
        return nil, errors.New("account not found")
    }
    if account.UserID != userID {
        return nil, errors.New("access denied")
    }
    return account, nil
}

func (s *AccountService) Deposit(userID int, req *models.DepositRequest) (*models.Transaction, error) {
    account, err := s.GetAccount(userID, req.AccountID)
    if err != nil {
        return nil, err
    }
    
    if err := s.accountRepo.UpdateBalance(account.ID, req.Amount); err != nil {
        return nil, err
    }
    
    tx := &models.Transaction{
        ToAccountID: &account.ID,
        Amount:      req.Amount,
        Type:        "deposit",
        Description: "Deposit to account",
    }
    
    if err := s.transactionRepo.Create(tx); err != nil {
        return nil, err
    }
    
    return tx, nil
}

func (s *AccountService) Transfer(userID int, req *models.TransferRequest) (*models.Transaction, error) {
    // Проверяем отправителя
    fromAccount, err := s.GetAccount(userID, req.FromAccountID)
    if err != nil {
        return nil, err
    }
    
    // Проверяем получателя
    toAccount, err := s.accountRepo.FindByID(req.ToAccountID)
    if err != nil {
        return nil, err
    }
    if toAccount == nil {
        return nil, errors.New("recipient account not found")
    }
    
    // Проверяем баланс
    if fromAccount.Balance < req.Amount {
        return nil, errors.New("insufficient funds")
    }
    
    // Списываем со счета отправителя
    if err := s.accountRepo.UpdateBalance(fromAccount.ID, -req.Amount); err != nil {
        return nil, err
    }
    
    // Зачисляем на счет получателя
    if err := s.accountRepo.UpdateBalance(toAccount.ID, req.Amount); err != nil {
        // Откатываем списание
        s.accountRepo.UpdateBalance(fromAccount.ID, req.Amount)
        return nil, err
    }
    
    tx := &models.Transaction{
        FromAccountID: &fromAccount.ID,
        ToAccountID:   &toAccount.ID,
        Amount:        req.Amount,
        Type:          "transfer",
        Description:   req.Description,
    }
    
    if err := s.transactionRepo.Create(tx); err != nil {
        return nil, err
    }
    
    return tx, nil
}

func (s *AccountService) generateAccountNumber() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("40817810%010d", rand.Intn(10000000000))
}
