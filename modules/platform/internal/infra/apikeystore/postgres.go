package apikeystore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jeielsantos/aureon/modules/platform/internal/domain"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Save(ctx context.Context, key domain.APIKey) error {
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO api_keys (id, account_id, secret_hash, created_at, revoked_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET revoked_at = EXCLUDED.revoked_at
	`, string(key.ID), string(key.AccountID), key.SecretHash, key.CreatedAt, key.RevokedAt)
	if err != nil {
		return fmt.Errorf("apikeystore: save key %s: %w", key.ID, err)
	}
	return nil
}

func (p *Postgres) FindBySecretHash(ctx context.Context, secretHash string) (domain.APIKey, error) {
	var (
		id, accountID string
		key           domain.APIKey
		revokedAt     sql.NullTime
	)

	err := p.db.QueryRowContext(ctx, `
		SELECT id, account_id, secret_hash, created_at, revoked_at
		FROM api_keys WHERE secret_hash = $1
	`, secretHash).Scan(&id, &accountID, &key.SecretHash, &key.CreatedAt, &revokedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.APIKey{}, fmt.Errorf("apikeystore: no matching api key")
	}
	if err != nil {
		return domain.APIKey{}, fmt.Errorf("apikeystore: find by secret hash: %w", err)
	}

	key.ID = domain.APIKeyID(id)
	key.AccountID = domain.AccountID(accountID)
	if revokedAt.Valid {
		key.RevokedAt = &revokedAt.Time
	}

	return key, nil
}

func (p *Postgres) Revoke(ctx context.Context, id domain.APIKeyID) error {
	res, err := p.db.ExecContext(ctx, `
		UPDATE api_keys SET revoked_at = now() WHERE id = $1 AND revoked_at IS NULL
	`, string(id))
	if err != nil {
		return fmt.Errorf("apikeystore: revoke key %s: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("apikeystore: revoke key %s: %w", id, err)
	}
	if n == 0 {
		return fmt.Errorf("apikeystore: no active api key with id %q", id)
	}

	return nil
}
