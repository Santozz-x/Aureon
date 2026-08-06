package domain

import "time"

type AccountID string

type Account struct {
	ID        AccountID
	Name      string
	CreatedAt time.Time
}

type APIKeyID string

// APIKey never holds the raw secret — only its hash. The raw secret is
// generated once by usecase.APIKeyService.Issue and shown to the caller
// exactly once, the same pattern GitHub/Stripe use for API tokens.
type APIKey struct {
	ID         APIKeyID
	AccountID  AccountID
	SecretHash string
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

func (k APIKey) Revoked() bool {
	return k.RevokedAt != nil
}
