package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	SMTPServer = "125.235.240.36" // Server SMTP của Viettel
	SMTPPort   = "465"            // Cổng SMTP (có thể đổi thành 465 nếu dùng SSL)
)

type EmailConfig struct {
	Username string   `json:"username"`
	ToEmails []string `json:"to_emails"`
	Password string   `json:"password"`
}

type DataLicense struct {
	Api   string
	Value float64
}

var (
	Username string
	ToEmail  []string
	Password string
)

func LoadEmailConfig() error {
	path := os.Getenv("EMAIL_CONFIG_PATH")
	if path == "" {
		path = "/app/config/email_config.json"
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read email config file: %w", err)
	}
	var config EmailConfig
	if err := json.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("failed to parse email config: %w", err)
	}

	Username = config.Username
	ToEmail = config.ToEmails
	Password = config.Password

	return nil
}

type Mail struct {
	Subject string
	Body    string
}
