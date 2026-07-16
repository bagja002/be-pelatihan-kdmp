package config

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from the environment.
type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret       string
	JWTAccessTTL    time.Duration
	JWTRefreshTTL   time.Duration
	PesertaTokenTTL time.Duration

	CORSAllowedOrigins string

	RateLimitMax    int
	RateLimitWindow time.Duration

	BodyLimit int

	// Upload berkas (CV & sertifikat pelatih).
	UploadDir      string
	MaxUploadBytes int64

	// Upload berkas bahan ajar (PDF & PPT — PPT biasanya besar).
	MaxBahanAjarBytes int64

	// EncryptionKey is the 32-byte AES-256 key used for field-level encryption.
	EncryptionKey []byte
}

// Insecure development defaults. These are ONLY used when APP_ENV is not
// "production"; in production the corresponding env vars are mandatory.
const (
	devJWTSecret = "dev-insecure-secret-change-me-please-0123456789"
	devEncKeyHex = "0000000000000000000000000000000000000000000000000000000000000000"
)

// Load reads .env (if present) and the environment, applies defaults, and
// validates security-critical secrets. It fails fast in production when a
// required secret is missing or weak.
func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "3000"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "data-knmp"),

		JWTSecret:          getEnv("JWT_SECRET", ""),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}

	c.JWTAccessTTL = getDuration("JWT_ACCESS_TTL", 15*time.Minute)
	c.JWTRefreshTTL = getDuration("JWT_REFRESH_TTL", 7*24*time.Hour)
	c.PesertaTokenTTL = getDuration("PESERTA_TOKEN_TTL", 30*time.Minute)
	c.RateLimitMax = getInt("RATE_LIMIT_MAX", 100)
	c.RateLimitWindow = getDuration("RATE_LIMIT_WINDOW", time.Minute)
	c.BodyLimit = getInt("BODY_LIMIT_BYTES", 70*1024*1024) // 70 MiB (muat PDF + PPT bahan ajar)
	c.UploadDir = getEnv("UPLOAD_DIR", "./uploads")
	c.MaxUploadBytes = int64(getInt("MAX_UPLOAD_BYTES", 5*1024*1024))          // 5 MiB per berkas
	c.MaxBahanAjarBytes = int64(getInt("MAX_BAHAN_AJAR_BYTES", 30*1024*1024)) // 30 MiB per berkas bahan ajar

	if err := c.resolveSecrets(); err != nil {
		return nil, err
	}
	return c, nil
}

// IsProduction reports whether the app runs in production mode.
func (c *Config) IsProduction() bool { return c.AppEnv == "production" }

func (c *Config) resolveSecrets() error {
	prod := c.IsProduction()

	// --- JWT secret ---
	if c.JWTSecret == "" {
		if prod {
			return fmt.Errorf("JWT_SECRET is required in production")
		}
		log.Println("[SECURITY WARNING] JWT_SECRET not set — using insecure development default")
		c.JWTSecret = devJWTSecret
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	// --- Encryption key ---
	rawKey := getEnv("ENCRYPTION_KEY", "")
	if rawKey == "" {
		if prod {
			return fmt.Errorf("ENCRYPTION_KEY is required in production")
		}
		log.Println("[SECURITY WARNING] ENCRYPTION_KEY not set — using insecure development default")
		rawKey = devEncKeyHex
	}
	key, err := parseKey(rawKey)
	if err != nil {
		return fmt.Errorf("invalid ENCRYPTION_KEY: %w", err)
	}
	c.EncryptionKey = key
	return nil
}

// parseKey accepts a 32-byte key encoded as hex (64 chars) or standard base64.
func parseKey(s string) ([]byte, error) {
	if len(s) == 64 {
		if b, err := hex.DecodeString(s); err == nil && len(b) == 32 {
			return b, nil
		}
	}
	if b, err := base64.StdEncoding.DecodeString(s); err == nil && len(b) == 32 {
		return b, nil
	}
	return nil, fmt.Errorf("key must decode to exactly 32 bytes (hex of 64 chars or standard base64)")
}

// DSN builds the MySQL data source name for GORM.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("[config] invalid int for %s=%q, using default %d", key, v, fallback)
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		log.Printf("[config] invalid duration for %s=%q, using default %s", key, v, fallback)
	}
	return fallback
}
