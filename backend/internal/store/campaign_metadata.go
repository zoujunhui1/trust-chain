// campaign_metadata.go — the second deliberate exception alongside users.go
// (see that file's comment): campaign title/description/theme aren't
// rebuilt from on-chain events either, so this is the API's write path for
// them, kept out of store.go/read.go's indexer-writes / API-reads split.
package store

import (
	"context"
	"database/sql"
	"fmt"
)

// SetCampaignMetadata stores (or replaces) a campaign's title/description/
// theme. description and theme are optional — pass nil to leave them unset.
// Deliberately not gated on the campaigns row existing yet — see
// db/schema.sql's comment on campaign_metadata for why there's no FK.
func (s *Store) SetCampaignMetadata(ctx context.Context, campaignID uint64, title string, description, theme *string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO campaign_metadata (campaign_id, title, description, theme) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE title = VALUES(title), description = VALUES(description), theme = VALUES(theme)`,
		campaignID, title, nullableString(description), nullableString(theme))
	if err != nil {
		return fmt.Errorf("store: 写活动元数据失败 / set campaign metadata: %w", err)
	}
	return nil
}

func nullableString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
