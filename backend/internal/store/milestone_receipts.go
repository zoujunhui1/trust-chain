// milestone_receipts.go — a third table in the campaign_metadata.go family:
// not rebuilt from on-chain events, kept out of store.go/read.go's
// indexer-writes / API-reads split. See db/schema.sql's comment on
// milestone_receipts for why this exists (the on-chain receiptHash is a
// bytes32 fingerprint, never the actual file).
package store

import (
	"context"
	"database/sql"
	"fmt"
)

// MilestoneReceiptFile stores where a milestone's uploaded expense document
// lives on disk, alongside its backend-computed hash.
type MilestoneReceiptFile struct {
	FilePath    string
	FileName    string
	ContentType string
	FileSize    int64
	SHA256      string // hex, no 0x prefix
	Note        *string
}

// SetMilestoneReceiptFile records (or replaces, if re-uploaded) the file
// backing a milestone's receipt. Deliberately not gated on the milestone row
// existing yet, matching every other off-chain table's LEFT JOIN read.
func (s *Store) SetMilestoneReceiptFile(ctx context.Context, campaignID uint64, idx uint32, f MilestoneReceiptFile) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO milestone_receipts (campaign_id, milestone_idx, file_path, file_name, content_type, file_size, sha256, note)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			file_path = VALUES(file_path), file_name = VALUES(file_name), content_type = VALUES(content_type),
			file_size = VALUES(file_size), sha256 = VALUES(sha256), note = VALUES(note)`,
		campaignID, idx, f.FilePath, f.FileName, f.ContentType, f.FileSize, f.SHA256, nullableString(f.Note))
	if err != nil {
		return fmt.Errorf("store: 写里程碑凭证失败 / set milestone receipt file: %w", err)
	}
	return nil
}

// GetMilestoneReceiptFile looks up one milestone's stored receipt file, if
// any — used by the download handler to resolve a disk path from campaign+idx.
func (s *Store) GetMilestoneReceiptFile(ctx context.Context, campaignID uint64, idx uint32) (MilestoneReceiptFile, bool, error) {
	var f MilestoneReceiptFile
	var note sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT file_path, file_name, content_type, file_size, sha256, note
		FROM milestone_receipts WHERE campaign_id = ? AND milestone_idx = ?`,
		campaignID, idx,
	).Scan(&f.FilePath, &f.FileName, &f.ContentType, &f.FileSize, &f.SHA256, &note)
	switch err {
	case nil:
		if note.Valid {
			f.Note = &note.String
		}
		return f, true, nil
	case sql.ErrNoRows:
		return MilestoneReceiptFile{}, false, nil
	default:
		return MilestoneReceiptFile{}, false, fmt.Errorf("store: 查里程碑凭证失败 / get milestone receipt file: %w", err)
	}
}
