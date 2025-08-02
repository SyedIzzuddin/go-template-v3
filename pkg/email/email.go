package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds email service configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Service represents email service interface
type Service interface {
	SendVerificationEmail(toEmail, toName, verificationToken string) error
	SendPasswordResetEmail(toEmail, toName, resetToken string) error
}

// SMTPService implements email service using SMTP
type SMTPService struct {
	config *Config
}

// NewSMTPService creates a new SMTP email service
func NewSMTPService(config *Config) *SMTPService {
	return &SMTPService{
		config: config,
	}
}

// SendVerificationEmail sends an email verification email
func (s *SMTPService) SendVerificationEmail(toEmail, toName, verificationToken string) error {
	subject := "Verify Your Email Address"
	
	// Generate verification URL (this should come from config in real implementation)
	verificationURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?token=%s", verificationToken)
	
	body := s.generateVerificationEmailBody(toName, verificationURL)
	
	return s.sendEmail(toEmail, subject, body)
}

// SendPasswordResetEmail sends a password reset email
func (s *SMTPService) SendPasswordResetEmail(toEmail, toName, resetToken string) error {
	subject := "Reset Your Password"
	
	// Generate password reset URL (this should come from config in real implementation)
	resetURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/reset-password?token=%s", resetToken)
	
	body := s.generatePasswordResetEmailBody(toName, resetURL)
	
	return s.sendEmail(toEmail, subject, body)
}

// sendEmail sends an email using SMTP
func (s *SMTPService) sendEmail(to, subject, body string) error {
	// Create authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	
	// Create the email message
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	msg := s.buildEmailMessage(from, to, subject, body)
	
	// Send the email
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(msg))
	
	if err != nil {
		return fmt.Errorf("failed to send email to %s via %s: %w", to, addr, err)
	}
	
	return nil
}

// buildEmailMessage builds the complete email message
func (s *SMTPService) buildEmailMessage(from, to, subject, body string) string {
	var msg strings.Builder
	
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	
	return msg.String()
}

// generateVerificationEmailBody generates HTML email body for verification
func (s *SMTPService) generateVerificationEmailBody(name, verificationURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Go Template!</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Thank you for registering with Go Template! To complete your registration, please verify your email address by clicking the button below:</p>
            
            <a href="%s" class="button">Verify My Email</a>
            
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p><a href="%s">%s</a></p>
            
            <p>This verification link will expire in 24 hours for security reasons.</p>
            
            <p>If you didn't create an account with us, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2025 Go Template. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, name, verificationURL, verificationURL, verificationURL)
}

// generatePasswordResetEmailBody generates HTML email body for password reset
func (s *SMTPService) generatePasswordResetEmailBody(name, resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #FF6B6B; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #FF6B6B; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { padding: 20px; text-align: center; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>We received a request to reset your password. If you made this request, click the button below to reset your password:</p>
            
            <a href="%s" class="button">Reset My Password</a>
            
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p><a href="%s">%s</a></p>
            
            <p>This password reset link will expire in 24 hours for security reasons.</p>
            
            <p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>© 2025 Go Template. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, name, resetURL, resetURL, resetURL)
}