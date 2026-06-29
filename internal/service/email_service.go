package service

import (
    "crypto/tls"
    "fmt"
    "log"
    "strconv"
    
    "github.com/go-mail/mail/v2"
    "bank-service/internal/config"
)

type EmailService struct {
    config *config.Config
}

func NewEmailService(config *config.Config) *EmailService {
    return &EmailService{config: config}
}

func (s *EmailService) SendPaymentConfirmation(userEmail string, amount float64, description string) error {
    subject := "Payment Confirmation"
    body := fmt.Sprintf(`
        <h1>Payment Successful</h1>
        <p>Amount: <strong>%.2f RUB</strong></p>
        <p>Description: %s</p>
        <small>This is an automated notification</small>
    `, amount, description)
    
    return s.sendEmail(userEmail, subject, body)
}

func (s *EmailService) SendWelcomeEmail(userEmail, username string) error {
    subject := "Welcome to Bank Service!"
    body := fmt.Sprintf(`
        <h1>Welcome %s!</h1>
        <p>Your account has been successfully created.</p>
        <p>You can now use our banking services.</p>
        <small>This is an automated notification</small>
    `, username)
    
    return s.sendEmail(userEmail, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
    if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
        log.Println("SMTP not configured, skipping email")
        return nil
    }
    
    // Convert port to int
    port, err := strconv.Atoi(s.config.SMTPPort)
    if err != nil {
        return fmt.Errorf("invalid SMTP port: %v", err)
    }
    
    m := mail.NewMessage()
    m.SetHeader("From", s.config.SMTPUser)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)
    
    d := mail.NewDialer(
        s.config.SMTPHost,
        port,
        s.config.SMTPUser,
        s.config.SMTPPassword,
    )
    
    d.TLSConfig = &tls.Config{
        ServerName:         s.config.SMTPHost,
        InsecureSkipVerify: false,
    }
    
    if err := d.DialAndSend(m); err != nil {
        log.Printf("SMTP error: %v", err)
        return fmt.Errorf("email sending failed: %w", err)
    }
    
    log.Printf("Email sent to %s", to)
    return nil
}
