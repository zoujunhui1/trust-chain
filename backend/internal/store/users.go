// users.go — the one deliberate exception to store.go/read.go's write/read
// split (see their doc comments): `users` isn't rebuilt from on-chain events,
// it's written directly by the API when a wallet connects. Kept in its own
// file so that exception doesn't get lost among the indexer-only writes in
// store.go.
package store

import (
	"context"
	"fmt"
)

// UpsertUser records that `addr` connected with the given role snapshot —
// inserts a first-seen row, or just bumps last_connected_at (and refreshes
// role) on a repeat connection. Never a source of truth for permissions;
// see the users table's comment in db/schema.sql.
func (s *Store) UpsertUser(ctx context.Context, addr, role string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO users (address, role) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role = VALUES(role)`, addr, role)
	if err != nil {
		return fmt.Errorf("store: 写用户失败 / upsert user: %w", err)
	}
	return nil
}
