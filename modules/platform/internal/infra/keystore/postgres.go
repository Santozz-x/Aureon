package keystore

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"

	chainport "github.com/jeielsantos/aureon/modules/contracts"
)

// Postgres is a chainport.KeyStore backed by PostgreSQL. Private key
// material is encrypted with AES-256-GCM before it ever reaches the
// database — the wallet_keys.encrypted_private_key column never holds
// plaintext key bytes. The AES key itself currently comes straight from
// AUREON_KEYSTORE_ENCRYPTION_KEY (an env var), not a KMS/HSM. See TR-010
// in docs/tradeoffs.md for why that's still not sufficient for a
// production deployment holding real funds.
type Postgres struct {
	db     *sql.DB
	cipher cipher.AEAD
}

func NewPostgres(db *sql.DB, encryptionKey []byte) (*Postgres, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keystore: init cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keystore: init gcm: %w", err)
	}

	return &Postgres{db: db, cipher: gcm}, nil
}

func (p *Postgres) Put(ctx context.Context, address chainport.Address, privateKey []byte) error {
	ciphertext, err := p.encrypt(privateKey)
	if err != nil {
		return fmt.Errorf("keystore: encrypt key for %s: %w", address, err)
	}

	_, err = p.db.ExecContext(ctx, `
		INSERT INTO wallet_keys (address, encrypted_private_key)
		VALUES ($1, $2)
		ON CONFLICT (address) DO UPDATE SET encrypted_private_key = EXCLUDED.encrypted_private_key
	`, string(address), ciphertext)
	if err != nil {
		return fmt.Errorf("keystore: store key for %s: %w", address, err)
	}

	return nil
}

func (p *Postgres) Get(ctx context.Context, address chainport.Address) ([]byte, error) {
	var ciphertext []byte
	err := p.db.QueryRowContext(ctx, `
		SELECT encrypted_private_key FROM wallet_keys WHERE address = $1
	`, string(address)).Scan(&ciphertext)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("keystore: no key stored for %s", address)
	}
	if err != nil {
		return nil, fmt.Errorf("keystore: get key for %s: %w", address, err)
	}

	plaintext, err := p.decrypt(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("keystore: decrypt key for %s: %w", address, err)
	}

	return plaintext, nil
}

func (p *Postgres) encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, p.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return p.cipher.Seal(nonce, nonce, plaintext, nil), nil
}

func (p *Postgres) decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := p.cipher.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext shorter than nonce size")
	}
	nonce, data := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return p.cipher.Open(nil, nonce, data, nil)
}
