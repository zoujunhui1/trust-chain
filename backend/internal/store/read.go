// read.go 是 store 的「读」侧：给 REST API 用的查询方法。
// read.go is the "read" side of store: query methods used by the REST API.
//
// 与写侧的分工 / split from the write side:
//   - 写方法（store.go）只被索引器调用，把链上事件落库。
//   - 读方法（本文件）只被 API 调用，把库里的数据查出来给前端。
//   这正是「读写分离」在代码层面的映射。
//   This mirrors the read/write split at the code level.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// --- 读返回的视图结构（金额一律用 string，避免 JSON number 溢出）---
// --- read view structs (amounts as string to avoid JSON number overflow) ---

// CampaignView 是一个活动的对外视图。/ CampaignView is a campaign's outward view.
type CampaignView struct {
	ID             uint64    `json:"id"`
	Charity        string    `json:"charity"`
	Goal           string    `json:"goal"`     // wei，十进制字符串 / wei as decimal string
	Raised         string    `json:"raised"`   // wei
	Released       string    `json:"released"` // wei
	MilestoneCount uint32    `json:"milestoneCount"`
	MetadataHash   string    `json:"metadataHash"`
	Completed      bool      `json:"completed"`
	CreatedBlock   uint64    `json:"createdBlock"`
	CreatedTx      string    `json:"createdTx"`
	CreatedAt      time.Time `json:"createdAt"`
}

// MilestoneView 是一个里程碑的对外视图。/ MilestoneView is a milestone's outward view.
type MilestoneView struct {
	Idx         uint32  `json:"idx"`
	Amount      string  `json:"amount"` // wei
	State       uint8   `json:"state"`  // 0=Locked 1=Released 2=Proven
	ReceiptHash *string `json:"receiptHash"`
}

// DonationView 是一笔捐款的对外视图。/ DonationView is a single donation's outward view.
type DonationView struct {
	Donor       string    `json:"donor"`
	Amount      string    `json:"amount"` // wei
	BlockNumber uint64    `json:"blockNumber"`
	TxHash      string    `json:"txHash"`
	LogIndex    uint32    `json:"logIndex"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CharityView 是机构准入状态的对外视图。/ CharityView is a charity's verification status.
type CharityView struct {
	Address      string `json:"address"`
	Verified     bool   `json:"verified"`
	UpdatedBlock uint64 `json:"updatedBlock"`
}

// ListCampaigns 返回活动列表，按创建区块倒序（最新在前）。
// ListCampaigns returns campaigns, newest first (by creation block).
func (s *Store) ListCampaigns(ctx context.Context, limit, offset int) ([]CampaignView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, charity, goal, raised, released, milestone_count,
		       COALESCE(metadata_hash, ''), completed, created_block, created_tx, created_at
		FROM campaigns
		ORDER BY created_block DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("store: 查活动列表失败 / list campaigns: %w", err)
	}
	defer rows.Close()

	// 初始化成非 nil 空切片，JSON 序列化时是 []（而不是 null）。
	// Init as non-nil empty slice so JSON marshals to [] (not null).
	list := []CampaignView{}
	for rows.Next() {
		var c CampaignView
		if err := rows.Scan(
			&c.ID, &c.Charity, &c.Goal, &c.Raised, &c.Released, &c.MilestoneCount,
			&c.MetadataHash, &c.Completed, &c.CreatedBlock, &c.CreatedTx, &c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: 扫描活动失败 / scan campaign: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// GetCampaign 返回单个活动；found=false 表示不存在（供 API 回 404）。
// GetCampaign returns one campaign; found=false means not found (API returns 404).
func (s *Store) GetCampaign(ctx context.Context, id uint64) (c CampaignView, found bool, err error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, charity, goal, raised, released, milestone_count,
		       COALESCE(metadata_hash, ''), completed, created_block, created_tx, created_at
		FROM campaigns WHERE id = ?`, id)
	switch err = row.Scan(
		&c.ID, &c.Charity, &c.Goal, &c.Raised, &c.Released, &c.MilestoneCount,
		&c.MetadataHash, &c.Completed, &c.CreatedBlock, &c.CreatedTx, &c.CreatedAt,
	); err {
	case nil:
		return c, true, nil
	case sql.ErrNoRows:
		return CampaignView{}, false, nil
	default:
		return CampaignView{}, false, fmt.Errorf("store: 查活动失败 / get campaign: %w", err)
	}
}

// ListMilestones 返回某活动的里程碑，按序号升序。
// ListMilestones returns a campaign's milestones ordered by index.
func (s *Store) ListMilestones(ctx context.Context, campaignID uint64) ([]MilestoneView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT idx, amount, state, receipt_hash
		FROM milestones WHERE campaign_id = ?
		ORDER BY idx ASC`, campaignID)
	if err != nil {
		return nil, fmt.Errorf("store: 查里程碑失败 / list milestones: %w", err)
	}
	defer rows.Close()

	list := []MilestoneView{}
	for rows.Next() {
		var m MilestoneView
		// receipt_hash 可能为 NULL，用 sql.NullString 接。
		// receipt_hash may be NULL; scan into sql.NullString.
		var receipt sql.NullString
		if err := rows.Scan(&m.Idx, &m.Amount, &m.State, &receipt); err != nil {
			return nil, fmt.Errorf("store: 扫描里程碑失败 / scan milestone: %w", err)
		}
		if receipt.Valid {
			m.ReceiptHash = &receipt.String
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// ListDonations 返回某活动的捐款流水，按区块+日志序号倒序（最新在前）。
// ListDonations returns a campaign's donations, newest first.
func (s *Store) ListDonations(ctx context.Context, campaignID uint64, limit, offset int) ([]DonationView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT donor, amount, block_number, tx_hash, log_index, created_at
		FROM donations WHERE campaign_id = ?
		ORDER BY block_number DESC, log_index DESC
		LIMIT ? OFFSET ?`, campaignID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("store: 查捐款失败 / list donations: %w", err)
	}
	defer rows.Close()

	list := []DonationView{}
	for rows.Next() {
		var d DonationView
		if err := rows.Scan(&d.Donor, &d.Amount, &d.BlockNumber, &d.TxHash, &d.LogIndex, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("store: 扫描捐款失败 / scan donation: %w", err)
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

// GetCharity 返回机构准入状态；found=false 表示查无此机构（从未认证过）。
// GetCharity returns a charity's status; found=false means unknown (never seen).
func (s *Store) GetCharity(ctx context.Context, addr string) (c CharityView, found bool, err error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT address, verified, updated_block FROM charities WHERE address = ?`, addr)
	switch err = row.Scan(&c.Address, &c.Verified, &c.UpdatedBlock); err {
	case nil:
		return c, true, nil
	case sql.ErrNoRows:
		return CharityView{}, false, nil
	default:
		return CharityView{}, false, fmt.Errorf("store: 查机构失败 / get charity: %w", err)
	}
}
