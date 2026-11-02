package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL   string
	HTTPAddr      string
	SessionSecret string
	CookieSecure  bool
	WebRoot       string
	S3Endpoint       string
	S3PublicEndpoint string
	S3Region         string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	S3PathStyle      bool
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	SMTPFrom      string
}

func Load() (Config, error) {
	loadDotEnv(".env")
	cfg := Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		CookieSecure:  os.Getenv("COOKIE_SECURE") == "true",
		WebRoot:       env("WEB_ROOT", "web/dist"),
		S3Endpoint:       os.Getenv("S3_ENDPOINT"),
		S3PublicEndpoint: os.Getenv("S3_PUBLIC_ENDPOINT"),
		S3Region:         env("S3_REGION", "us-east-1"),
		S3Bucket:      os.Getenv("S3_BUCKET"),
		S3AccessKey:   os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("S3_SECRET_KEY"),
		S3PathStyle:   os.Getenv("S3_PATH_STYLE") == "true",
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      env("SMTP_PORT", "587"),
		SMTPUser:      os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:      env("SMTP_FROM", "diary@localhost"),
	}
	if cfg.DatabaseURL == "" || cfg.SessionSecret == "" || cfg.S3Bucket == "" || cfg.S3AccessKey == "" || cfg.S3SecretKey == "" {
		return Config{}, fmt.Errorf("DATABASE_URL, SESSION_SECRET, S3_BUCKET, S3_ACCESS_KEY, and S3_SECRET_KEY are required")
	}
	if len(cfg.SessionSecret) < 16 {
		return Config{}, fmt.Errorf("SESSION_SECRET must be at least 16 characters")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		os.Setenv(key, strings.TrimSpace(val))
	}
}
