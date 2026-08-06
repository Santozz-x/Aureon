package config

import "testing"

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("AUREON_DATABASE_URL", "")
	t.Setenv("AUREON_KEYSTORE_ENCRYPTION_KEY", "db29b037c63f986365949818a2959a9e3002234f39d972dfe1c7cca6b9f46b16")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing AUREON_DATABASE_URL, got nil")
	}
}

func TestLoad_MissingEncryptionKey(t *testing.T) {
	t.Setenv("AUREON_DATABASE_URL", "postgres://user:pass@localhost:5432/aureon")
	t.Setenv("AUREON_KEYSTORE_ENCRYPTION_KEY", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing AUREON_KEYSTORE_ENCRYPTION_KEY, got nil")
	}
}

func TestLoad_EncryptionKeyWrongLength(t *testing.T) {
	t.Setenv("AUREON_DATABASE_URL", "postgres://user:pass@localhost:5432/aureon")
	t.Setenv("AUREON_KEYSTORE_ENCRYPTION_KEY", "aabbcc")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for short AUREON_KEYSTORE_ENCRYPTION_KEY, got nil")
	}
}

func TestLoad_EncryptionKeyNotHex(t *testing.T) {
	t.Setenv("AUREON_DATABASE_URL", "postgres://user:pass@localhost:5432/aureon")
	t.Setenv("AUREON_KEYSTORE_ENCRYPTION_KEY", "not-hex-zzzz")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for non-hex AUREON_KEYSTORE_ENCRYPTION_KEY, got nil")
	}
}

func TestLoad_Valid(t *testing.T) {
	t.Setenv("AUREON_DATABASE_URL", "postgres://user:pass@localhost:5432/aureon")
	t.Setenv("AUREON_KEYSTORE_ENCRYPTION_KEY", "db29b037c63f986365949818a2959a9e3002234f39d972dfe1c7cca6b9f46b16")
	t.Setenv("AUREON_HTTP_ADDR", "")
	t.Setenv("AUREON_ARC_RPC_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %s, want default :8080", cfg.HTTPAddr)
	}
	if cfg.ArcRPCURL != defaultArcRPCURL {
		t.Fatalf("ArcRPCURL = %s, want default %s", cfg.ArcRPCURL, defaultArcRPCURL)
	}
	if len(cfg.KeystoreEncryptionKey) != encryptionKeySize {
		t.Fatalf("KeystoreEncryptionKey length = %d, want %d", len(cfg.KeystoreEncryptionKey), encryptionKeySize)
	}
}
