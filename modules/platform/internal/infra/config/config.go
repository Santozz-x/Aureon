package config

import (
	"encoding/hex"
	"fmt"
	"os"
)

const (
	defaultArcRPCURL  = "https://rpc.testnet.arc.io"
	encryptionKeySize = 32 // AES-256
)

type Config struct {
	HTTPAddr              string
	ArcRPCURL             string
	DatabaseURL           string
	KeystoreEncryptionKey []byte
}

func Load() (Config, error) {
	databaseURL := os.Getenv("AUREON_DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("config: AUREON_DATABASE_URL is required")
	}

	encryptionKey, err := loadEncryptionKey()
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPAddr:              getenv("AUREON_HTTP_ADDR", ":8080"),
		ArcRPCURL:             getenv("AUREON_ARC_RPC_URL", defaultArcRPCURL),
		DatabaseURL:           databaseURL,
		KeystoreEncryptionKey: encryptionKey,
	}, nil
}

func loadEncryptionKey() ([]byte, error) {
	hexKey := os.Getenv("AUREON_KEYSTORE_ENCRYPTION_KEY")
	if hexKey == "" {
		return nil, fmt.Errorf("config: AUREON_KEYSTORE_ENCRYPTION_KEY is required (hex-encoded, %d bytes)", encryptionKeySize)
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("config: AUREON_KEYSTORE_ENCRYPTION_KEY must be hex-encoded: %w", err)
	}
	if len(key) != encryptionKeySize {
		return nil, fmt.Errorf("config: AUREON_KEYSTORE_ENCRYPTION_KEY must decode to %d bytes, got %d", encryptionKeySize, len(key))
	}

	return key, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
