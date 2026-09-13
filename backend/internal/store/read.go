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
	// Confirmed 是 false 时，这行是 API 乐观写入的（交易刚确认，indexer 还没
	// 追上），前端应该显示"正在确认"之类的提示。见 store.InsertCampaign。
	// False means this row was written optimistically by the API (tx just
	// confirmed, indexer hasn't caught up yet) — the frontend should show
	// something like "confirming on-chain…". See store.InsertCampaign.
	Confirmed bool `json:"confirmed"`
	// Title/Description/Theme 来自 campaign_metadata（LEFT JOIN），没设置过
	// 就是 nil —— 前端 Title 缺失时回退显示 "Campaign #{id}"，Theme 缺失时
	// 回退成按 id 确定性分配。见 store/campaign_metadata.go。
	// From campaign_metadata (LEFT JOIN); nil if never set — the frontend
	// falls back to "Campaign #{id}" for a missing title, and its old
	// deterministic-by-id pick for a missing theme. See
	// store/campaign_metadata.go.
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Theme       *string `json:"theme"`
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

// ActivityView 是全局活动流水的一行对外视图。
// ActivityView is one row of the global activity feed's outward view.
type ActivityView struct {
	ID           uint64    `json:"id"`
	CampaignID   uint64    `json:"campaignId"`
	EventType    string    `json:"eventType"`
	Amount       *string   `json:"amount"`       // wei; null for non-monetary events
	MilestoneIdx *uint32   `json:"milestoneIdx"` // null unless milestone-scoped
	BlockNumber  uint64    `json:"blockNumber"`
	TxHash       string    `json:"txHash"`
	LogIndex     uint32    `json:"logIndex"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ListCampaigns 返回活动列表，按创建区块倒序（最新在前）。
// ListCampaigns returns campaigns, newest first (by creation block).
func (s *Store) ListCampaigns(ctx context.Context, limit, offset int) ([]CampaignView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.charity, c.goal, c.raised, c.released, c.milestone_count,
		       COALESCE(c.metadata_hash, ''), c.completed, c.confirmed, c.created_block, c.created_tx, c.created_at,
		       m.title, m.description, m.theme
		FROM campaigns c
		LEFT JOIN campaign_metadata m ON m.campaign_id = c.id
		ORDER BY c.created_block DESC
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
		var title, description, theme sql.NullString
		if err := rows.Scan(
			&c.ID, &c.Charity, &c.Goal, &c.Raised, &c.Released, &c.MilestoneCount,
			&c.MetadataHash, &c.Completed, &c.Confirmed, &c.CreatedBlock, &c.CreatedTx, &c.CreatedAt,
			&title, &description, &theme,
		); err != nil {
			return nil, fmt.Errorf("store: 扫描活动失败 / scan campaign: %w", err)
		}
		if title.Valid {
			c.Title = &title.String
		}
		if description.Valid {
			c.Description = &description.String
		}
		if theme.Valid {
			c.Theme = &theme.String
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// GetCampaign 返回单个活动；found=false 表示不存在（供 API 回 404）。
// GetCampaign returns one campaign; found=false means not found (API returns 404).
func (s *Store) GetCampaign(ctx context.Context, id uint64) (c CampaignView, found bool, err error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT c.id, c.charity, c.goal, c.raised, c.released, c.milestone_count,
		       COALESCE(c.metadata_hash, ''), c.completed, c.confirmed, c.created_block, c.created_tx, c.created_at,
		       m.title, m.description, m.theme
		FROM campaigns c
		LEFT JOIN campaign_metadata m ON m.campaign_id = c.id
		WHERE c.id = ?`, id)
	var title, description, theme sql.NullString
	switch err = row.Scan(
		&c.ID, &c.Charity, &c.Goal, &c.Raised, &c.Released, &c.MilestoneCount,
		&c.MetadataHash, &c.Completed, &c.Confirmed, &c.CreatedBlock, &c.CreatedTx, &c.CreatedAt,
		&title, &description, &theme,
	); err {
	case nil:
		if title.Valid {
			c.Title = &title.String
		}
		if description.Valid {
			c.Description = &description.String
		}
		if theme.Valid {
			c.Theme = &theme.String
		}
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

// ListCharities 返回所有见过的机构（含已撤销的），按最后变更区块倒序。
// ListCharities returns every charity ever seen (including revoked ones), newest change first.
func (s *Store) ListCharities(ctx context.Context, limit, offset int) ([]CharityView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT address, verified, updated_block FROM charities
		ORDER BY updated_block DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("store: 查机构列表失败 / list charities: %w", err)
	}
	defer rows.Close()

	list := []CharityView{}
	for rows.Next() {
		var c CharityView
		if err := rows.Scan(&c.Address, &c.Verified, &c.UpdatedBlock); err != nil {
			return nil, fmt.Errorf("store: 扫描机构失败 / scan charity: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// ListActivity 返回全局活动流水，按区块+日志序号倒序（最新在前）。
// ListActivity returns the global activity feed, newest first.
func (s *Store) ListActivity(ctx context.Context, limit, offset int) ([]ActivityView, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, campaign_id, event_type, amount, milestone_idx, block_number, tx_hash, log_index, created_at
		FROM activity_events
		ORDER BY block_number DESC, log_index DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("store: 查活动流水失败 / list activity: %w", err)
	}
	defer rows.Close()

	list := []ActivityView{}
	for rows.Next() {
		var a ActivityView
		var amount sql.NullString
		var idx sql.NullInt64
		if err := rows.Scan(
			&a.ID, &a.CampaignID, &a.EventType, &amount, &idx, &a.BlockNumber, &a.TxHash, &a.LogIndex, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: 扫描活动流水失败 / scan activity: %w", err)
		}
		if amount.Valid {
			a.Amount = &amount.String
		}
		if idx.Valid {
			v := uint32(idx.Int64)
			a.MilestoneIdx = &v
		}
		list = append(list, a)
	}
	return list, rows.Err()
}
