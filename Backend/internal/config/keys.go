package config

import (
	"crypto/rand"

	"github.com/WBetEngine/Seosementara/Backend/internal/crypto"
)

// ResolveEncryptionKey prefers MASTER_ENCRYPTION_KEY then PIXEL_ENCRYPTION_KEY.
func (c Config) ResolveEncryptionKey() ([]byte, error) {
	if c.MasterEncryptionKey != "" {
		return crypto.KeyFromEnv(c.MasterEncryptionKey)
	}
	if c.PixelEncryptionKey != "" {
		return crypto.KeyFromEnv(c.PixelEncryptionKey)
	}
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		return nil, err
	}
	return k, nil
}
