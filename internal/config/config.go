package config

import (
    "os"
)

type Config struct {
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    JWTSecret      string
    PGPPrivateKey  string
    PGPPublicKey   string
    HMACSecret     string
    SMTPHost       string
    SMTPPort       string
    SMTPUser       string
    SMTPPassword   string
    ServerPort     string
    CentralBankURL string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", "password"),
        DBName:         getEnv("DB_NAME", "bankdb"),
        JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
        PGPPrivateKey:  getEnv("PGP_PRIVATE_KEY", ""),
        PGPPublicKey:   getEnv("PGP_PUBLIC_KEY", ""),
        HMACSecret:     getEnv("HMAC_SECRET", "hmac-secret-key"),
        SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
        SMTPPort:       getEnv("SMTP_PORT", "587"),
        SMTPUser:       getEnv("SMTP_USER", ""),
        SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
        ServerPort:     getEnv("SERVER_PORT", "8080"),
        CentralBankURL: getEnv("CENTRAL_BANK_URL", "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
