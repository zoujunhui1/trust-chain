// milestone_metadata.go — same family as campaign_metadata.go: not rebuilt
// from on-chain events, kept out of store.go/read.go's indexer-writes /
// API-reads split. See db/schema.sql's comment on milestone_metadata for
// why this exists and how it differs from milestone_receipts.note.
package store

import (
	"context"
	"fmt"
	"strings"
)

// SetMilestoneDescriptions stores each milestone's plan description,
// indexed by position (descriptions[i] is milestone i). Blank entries are
// skipped — this is an optional field, so there's nothing to store for a
// milestone the charity left empty.
func (s *Store) SetMilestoneDescriptions(ctx context.Context, campaignID uint64, descriptions []string) error {
	if len(descriptions) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: 开始事务失败 / begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO milestone_metadata (campaign_id, milestone_idx, description) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE description = VALUES(description)`)
	if err != nil {
		return fmt.Errorf("store: 准备语句失败 / prepare statement: %w", err)
	}
	defer stmt.Close()

	for idx, d := range descriptions {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		if _, err := stmt.ExecContext(ctx, campaignID, idx, d); err != nil {
			return fmt.Errorf("store: 写里程碑说明失败 / set milestone description: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: 提交事务失败 / commit tx: %w", err)
	}
	return nil
}
