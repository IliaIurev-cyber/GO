package main

import (
    "encoding/json"
    "net/http"
    "os"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    
    "bank-service/internal/config"
    "bank-service/internal/handler"
    "bank-service/internal/middleware"
    "bank-service/internal/repository"
    "bank-service/internal/service"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        logrus.Warn("No .env file found, using defaults")
    }
    
    // Load configuration
    cfg := config.LoadConfig()
    
    // Setup logging
    logrus.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    
    // Initialize repositories (in-memory)
    userRepo := repository.NewInMemoryUserRepository()
    accountRepo := repository.NewInMemoryAccountRepository()
    transactionRepo := repository.NewInMemoryTransactionRepository()
    logrus.Info("Using in-memory storage (no database required)")
    
    // Initialize services
    authService := service.NewAuthService(userRepo, cfg)
    accountService := service.NewAccountService(accountRepo, transactionRepo)
    cbrService := service.NewCBRService(cfg)
    emailService := service.NewEmailService(cfg)
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    accountHandler := handler.NewAccountHandler(accountService)
    
    // Create router
    r := mux.NewRouter()
    r.Use(middleware.LoggerMiddleware)
    
    // Public routes
    r.HandleFunc("/register", authHandler.Register).Methods("POST")
    r.HandleFunc("/login", authHandler.Login).Methods("POST")
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "ok",
            "time":   time.Now().Format(time.RFC3339),
        })
    }).Methods("GET")
    
    // Root endpoint
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "service": "Bank Service API",
            "version": "1.0.0",
            "status":  "running",
            "storage": "in-memory",
        })
    }).Methods("GET")
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{
            "message": "pong",
            "time":    time.Now().Format(time.RFC3339),
        })
    }).Methods("GET")
    
    // Central Bank rate endpoint (public)
    api.HandleFunc("/cbr-rate", func(w http.ResponseWriter, r *http.Request) {
        rate, err := cbrService.GetKeyRate()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]interface{}{
            "rate":     rate,
            "currency": "RUB",
            "date":     time.Now().Format(time.RFC3339),
        })
    }).Methods("GET")
    
    // Protected routes
    protected := api.PathPrefix("/").Subrouter()
    protected.Use(middleware.AuthMiddleware(authService))
    
    // Profile
    protected.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r.Context())
        json.NewEncoder(w).Encode(map[string]interface{}{
            "user_id": userID,
            "message": "Welcome to your profile! This is a protected endpoint.",
        })
    }).Methods("GET")
    
    // Account routes
    protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
    protected.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
    protected.HandleFunc("/accounts/{id:[0-9]+}", accountHandler.GetAccount).Methods("GET")
    protected.HandleFunc("/deposit", accountHandler.Deposit).Methods("POST")
    protected.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")
    
    // Send test email (protected)
    protected.HandleFunc("/test-email", func(w http.ResponseWriter, r *http.Request) {
        // In real app, get user email from database
        err := emailService.SendWelcomeEmail("test@example.com", "Test User")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Test email sent successfully",
        })
    }).Methods("POST")
    
    // Start server
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    
    logrus.Infof("Server starting on :%s", port)
    logrus.Infof("Health check: http://localhost:%s/health", port)
    logrus.Infof("API ping: http://localhost:%s/api/ping", port)
    logrus.Infof("CBR Rate: http://localhost:%s/api/cbr-rate", port)
    logrus.Infof("Register: http://localhost:%s/register", port)
    logrus.Infof("Login: http://localhost:%s/login", port)
    logrus.Infof("Accounts: http://localhost:%s/api/accounts", port)
    logrus.Infof("Deposit: http://localhost:%s/api/deposit", port)
    logrus.Infof("Transfer: http://localhost:%s/api/transfer", port)
    
    if err := http.ListenAndServe(":"+port, r); err != nil {
        logrus.Fatalf("Server failed: %v", err)
    }
}
