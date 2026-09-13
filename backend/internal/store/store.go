// Package store 是索引器的存储层：把解析好的链上事件幂等地写入 MySQL。
// Package store is the indexer's persistence layer: it writes parsed on-chain
// events into MySQL idempotently.
//
// 幂等如何保证 / how idempotency is guaranteed:
//   - donations 靠唯一键 (tx_hash, log_index) + INSERT IGNORE，重扫不会重复插入。
//   - 金额字段（raised/released）不做累加，而是用 SUM 重算，重放也得到同一结果。
//   - campaigns/milestones 用 INSERT IGNORE；状态推进用 GREATEST，只进不退。
package store

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	_ "github.com/go-sql-driver/mysql" // 注册 mysql driver / register the driver
)

// Store 封装一个数据库连接池。/ Store wraps a database connection pool.
type Store struct {
	db *sql.DB
}

// --- 传给写方法的数据结构（与链上事件字段对应）/ input structs mirroring events ---

type Campaign struct {
	ID             uint64
	Charity        string   // 小写 0x 地址
	Goal           *big.Int // wei
	MilestoneCount uint32
	MetadataHash   string // 0x + 64 hex（view 补齐）
	CreatedBlock   uint64
	CreatedTx      string
}

type Milestone struct {
	CampaignID uint64
	Idx        uint32
	Amount     *big.Int // wei
}

type Donation struct {
	CampaignID  uint64
	Donor       string
	Amount      *big.Int // wei
	BlockNumber uint64
	TxHash      string
	LogIndex    uint32
}

// ActivityEvent is one row of the global activity feed (see activity_events).
// Amount / MilestoneIdx are nil for event types that don't carry them.
type ActivityEvent struct {
	CampaignID   uint64
	EventType    string
	Amount       *big.Int
	MilestoneIdx *uint32
	BlockNumber  uint64
	TxHash       string
	LogIndex     uint32
}

// Open 建立连接池并 Ping 校验（fail fast）。/ Open dials the pool and pings (fail fast).
func Open(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: 打开数据库失败 / open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: 连接数据库失败 / ping db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return &Store{db: db}, nil
}

// Close 关闭连接池。/ Close releases the pool.
func (s *Store) Close() error { return s.db.Close() }

// --- 游标 / cursor ---

// GetCursor 返回该链已处理到的区块；found=false 表示还没有记录（首次运行）。
// GetCursor returns the last processed block; found=false means no record yet.
func (s *Store) GetCursor(ctx context.Context, id string) (block uint64, found bool, err error) {
	row := s.db.QueryRowContext(ctx, `SELECT last_block FROM indexer_cursor WHERE id = ?`, id)
	switch err = row.Scan(&block); err {
	case nil:
		return block, true, nil
	case sql.ErrNoRows:
		return 0, false, nil
	default:
		return 0, false, fmt.Errorf("store: 读游标失败 / get cursor: %w", err)
	}
}

// SaveCursor 更新（或插入）该链的处理进度。/ SaveCursor upserts the progress for a chain.
func (s *Store) SaveCursor(ctx context.Context, id string, block uint64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO indexer_cursor (id, last_block) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE last_block = VALUES(last_block)`, id, block)
	if err != nil {
		return fmt.Errorf("store: 存游标失败 / save cursor: %w", err)
	}
	return nil
}

// --- 慈善机构 / charities ---

// UpsertCharity 写入机构准入状态（CharityVerified / CharityRevoked）。
// UpsertCharity records a charity's verification status.
func (s *Store) UpsertCharity(ctx context.Context, addr string, verified bool, block uint64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO charities (address, verified, updated_block) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE verified = VALUES(verified), updated_block = VALUES(updated_block)`,
		addr, verified, block)
	if err != nil {
		return fmt.Errorf("store: 写机构失败 / upsert charity: %w", err)
	}
	return nil
}

// --- 活动 + 里程碑 / campaign + milestones ---

// InsertCampaign 在一个事务里插入/确认活动及其里程碑。
//
// confirmed 区分两种调用方 / confirmed distinguishes two callers:
//   - indexer 处理到真实链上事件时传 true：数据可信，会覆盖任何早先由 API
//     乐观写入（见 POST /api/campaigns）的字段，并把 confirmed 推成 true。
//   - API 乐观写入（交易刚确认，indexer 还没追上）传 false：如果这行已经被
//     indexer 确认过，这次调用不会覆盖/降级任何字段。
//   - When the indexer processes a real on-chain event it passes true: the
//     data is authoritative, overwrites any earlier optimistic write from
//     POST /api/campaigns, and pushes confirmed to true.
//   - The API's optimistic write (tx just confirmed, indexer hasn't caught
//     up yet) passes false: if the row is already indexer-confirmed, this
//     call never overwrites or downgrades any field.
//
// 里程碑本身没有这个真假区分——乐观写入提交的金额和最终链上事件的金额是
// 同一笔交易的数据，不存在"谁更可信"的问题，所以还是用 INSERT IGNORE。
// Milestones don't need this distinction — the amounts an optimistic write
// submits are from the very same transaction the eventual on-chain event
// describes, so there's nothing to reconcile; INSERT IGNORE is enough.
func (s *Store) InsertCampaign(ctx context.Context, c Campaign, ms []Milestone, confirmed bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: 开事务失败 / begin tx: %w", err)
	}
	defer tx.Rollback() // 已 Commit 后再 Rollback 是无操作 / no-op after commit

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO campaigns
			(id, charity, goal, milestone_count, metadata_hash, created_block, created_tx, confirmed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			charity         = IF(confirmed, charity, VALUES(charity)),
			goal            = IF(confirmed, goal, VALUES(goal)),
			milestone_count = IF(confirmed, milestone_count, VALUES(milestone_count)),
			metadata_hash   = IF(confirmed, metadata_hash, VALUES(metadata_hash)),
			created_block   = IF(confirmed, created_block, VALUES(created_block)),
			created_tx      = IF(confirmed, created_tx, VALUES(created_tx)),
			confirmed       = GREATEST(confirmed, VALUES(confirmed))`,
		c.ID, c.Charity, c.Goal.String(), c.MilestoneCount, c.MetadataHash, c.CreatedBlock, c.CreatedTx, confirmed,
	); err != nil {
		return fmt.Errorf("store: 插入活动失败 / insert campaign: %w", err)
	}

	for _, m := range ms {
		if _, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO milestones (campaign_id, idx, amount, state)
			VALUES (?, ?, ?, 0)`, // 初始 Locked=0
			m.CampaignID, m.Idx, m.Amount.String(),
		); err != nil {
			return fmt.Errorf("store: 插入里程碑失败 / insert milestone: %w", err)
		}
	}
	return tx.Commit()
}

// InsertDonation 幂等插入一笔捐款，并把该活动的 raised 重算为捐款总和。
// InsertDonation idempotently inserts a donation and recomputes the campaign's raised total.
func (s *Store) InsertDonation(ctx context.Context, d Donation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: 开事务失败 / begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
		INSERT IGNORE INTO donations
			(campaign_id, donor, amount, block_number, tx_hash, log_index)
		VALUES (?, ?, ?, ?, ?, ?)`,
		d.CampaignID, d.Donor, d.Amount.String(), d.BlockNumber, d.TxHash, d.LogIndex,
	); err != nil {
		return fmt.Errorf("store: 插入捐款失败 / insert donation: %w", err)
	}

	// raised = SUM(donations)：重算而非累加 → 即便重扫也不会多加。
	// Recompute rather than increment, so re-scans never double-count.
	if _, err = tx.ExecContext(ctx, `
		UPDATE campaigns
		SET raised = (SELECT COALESCE(SUM(amount), 0) FROM donations WHERE campaign_id = ?)
		WHERE id = ?`, d.CampaignID, d.CampaignID,
	); err != nil {
		return fmt.Errorf("store: 更新 raised 失败 / update raised: %w", err)
	}
	return tx.Commit()
}

// MarkMilestoneReleased 把里程碑推进到 Released，并重算 released 总额。
// MarkMilestoneReleased advances a milestone to Released and recomputes released total.
func (s *Store) MarkMilestoneReleased(ctx context.Context, campaignID uint64, idx uint32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: 开事务失败 / begin tx: %w", err)
	}
	defer tx.Rollback()

	// GREATEST(state,1)：只进不退，已是 Proven(2) 就保持 2 → 幂等、与顺序无关。
	// GREATEST keeps state monotonic; if already Proven(2) it stays 2.
	if _, err = tx.ExecContext(ctx, `
		UPDATE milestones SET state = GREATEST(state, 1)
		WHERE campaign_id = ? AND idx = ?`, campaignID, idx,
	); err != nil {
		return fmt.Errorf("store: 更新里程碑状态失败 / update milestone state: %w", err)
	}

	// released = 已放款(state>=1)里程碑金额之和。/ released = sum of amounts with state>=1.
	if _, err = tx.ExecContext(ctx, `
		UPDATE campaigns
		SET released = (SELECT COALESCE(SUM(amount), 0) FROM milestones WHERE campaign_id = ? AND state >= 1)
		WHERE id = ?`, campaignID, campaignID,
	); err != nil {
		return fmt.Errorf("store: 更新 released 失败 / update released: %w", err)
	}
	return tx.Commit()
}

// MarkMilestoneProven 把里程碑推进到 Proven 并记录收据哈希。
// MarkMilestoneProven advances a milestone to Proven and stores its receipt hash.
func (s *Store) MarkMilestoneProven(ctx context.Context, campaignID uint64, idx uint32, receiptHash string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE milestones SET state = 2, receipt_hash = ?
		WHERE campaign_id = ? AND idx = ?`, receiptHash, campaignID, idx)
	if err != nil {
		return fmt.Errorf("store: 记录收据失败 / mark proven: %w", err)
	}
	return nil
}

// MarkCampaignCompleted 标记活动完成（收到最后一个里程碑的收据）。
// MarkCampaignCompleted flags a campaign as complete (last milestone receipt received).
func (s *Store) MarkCampaignCompleted(ctx context.Context, campaignID uint64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE campaigns SET completed = TRUE WHERE id = ?`, campaignID)
	if err != nil {
		return fmt.Errorf("store: 标记完成失败 / mark completed: %w", err)
	}
	return nil
}

// --- 活动流水 / activity feed ---

// InsertActivity 幂等追加一行全局活动流水（INSERT IGNORE，唯一键同 donations）。
// InsertActivity idempotently appends one activity-feed row (INSERT IGNORE, same unique key pattern as donations).
func (s *Store) InsertActivity(ctx context.Context, e ActivityEvent) error {
	var amount sql.NullString
	if e.Amount != nil {
		amount = sql.NullString{String: e.Amount.String(), Valid: true}
	}
	var idx sql.NullInt64
	if e.MilestoneIdx != nil {
		idx = sql.NullInt64{Int64: int64(*e.MilestoneIdx), Valid: true}
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT IGNORE INTO activity_events
			(campaign_id, event_type, amount, milestone_idx, block_number, tx_hash, log_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.CampaignID, e.EventType, amount, idx, e.BlockNumber, e.TxHash, e.LogIndex,
	)
	if err != nil {
		return fmt.Errorf("store: 写活动流水失败 / insert activity: %w", err)
	}
	return nil
}
