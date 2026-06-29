package utils

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func ComputeHMAC(data string, secret []byte) string {
    h := hmac.New(sha256.New, secret)
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data, hmacValue string, secret []byte) bool {
    expectedHMAC := ComputeHMAC(data, secret)
    return hmac.Equal([]byte(hmacValue), []byte(expectedHMAC))
}
