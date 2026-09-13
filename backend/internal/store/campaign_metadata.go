// campaign_metadata.go — the second deliberate exception alongside users.go
// (see that file's comment): campaign titles aren't rebuilt from on-chain
// events either, so this is the API's write path for them, kept out of
// store.go/read.go's indexer-writes / API-reads split.
package store

import (
	"context"
	"fmt"
)

// SetCampaignTitle stores (or replaces) a campaign's title. Deliberately not
// gated on the campaigns row existing yet — see db/schema.sql's comment on
// campaign_metadata for why there's no FK.
func (s *Store) SetCampaignTitle(ctx context.Context, campaignID uint64, title string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO campaign_metadata (campaign_id, title) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE title = VALUES(title)`, campaignID, title)
	if err != nil {
		return fmt.Errorf("store: 写活动标题失败 / set campaign title: %w", err)
	}
	return nil
}
