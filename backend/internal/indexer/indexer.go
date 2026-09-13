// Package indexer 是索引器的核心：把链上日志扫下来、解析、落库、推进游标。
// Package indexer is the core loop: scan logs, parse them, persist, advance cursor.
//
// 它把前面各层接起来 / it wires the layers together:
//   config → 参数    chain → 拉日志/调 view    store → 幂等落库
package indexer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/zoujunhui1/trust-chain/backend/internal/chain"
	"github.com/zoujunhui1/trust-chain/backend/internal/config"
	"github.com/zoujunhui1/trust-chain/backend/internal/contracts"
	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

// Indexer 持有各依赖以及两个合约 ABI（用于按事件签名分发日志）。
// Indexer holds dependencies and the two ABIs (used to dispatch logs by event signature).
type Indexer struct {
	cfg         *config.Config
	chain       *chain.Client
	store       *store.Store
	escrowABI   *abi.ABI
	registryABI *abi.ABI
}

// New 构造 Indexer，并取出两个合约的 ABI 供分发使用。
// New builds an Indexer and loads both ABIs for dispatch.
func New(cfg *config.Config, ch *chain.Client, st *store.Store) (*Indexer, error) {
	escrowABI, err := contracts.TrustChainEscrowMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("indexer: 取 escrow ABI 失败 / get escrow abi: %w", err)
	}
	registryABI, err := contracts.CharityRegistryMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("indexer: 取 registry ABI 失败 / get registry abi: %w", err)
	}
	return &Indexer{cfg: cfg, chain: ch, store: st, escrowABI: escrowABI, registryABI: registryABI}, nil
}

// ============================ 上半：扫描循环 / the scan loop ============================

// Run 阻塞式运行主循环，直到 ctx 取消。/ Run blocks running the main loop until ctx is cancelled.
func (ix *Indexer) Run(ctx context.Context) error {
	from, err := ix.startBlock(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("[indexer] 起始区块 / start block = %d（确认数 %d，批大小 %d）\n",
		from, ix.cfg.Confirmations, ix.cfg.BatchSize)

	for {
		// 每轮先看是否被要求退出 / check for shutdown each iteration.
		if err := ctx.Err(); err != nil {
			return err
		}

		latest, err := ix.chain.LatestBlock(ctx)
		if err != nil {
			// 网络抖动等瞬时错误：等一下重试，不推进游标。
			// Transient error (e.g. network): back off and retry without advancing.
			fmt.Printf("[indexer] 取最新块失败，稍后重试 / latest block failed: %v\n", err)
			if err := ix.sleep(ctx); err != nil {
				return err
			}
			continue
		}

		// 落后链头 Confirmations 个块再处理，规避链重组。
		// Stay Confirmations blocks behind the head to avoid reorgs.
		if latest < ix.cfg.Confirmations {
			if err := ix.sleep(ctx); err != nil {
				return err
			}
			continue
		}
		safeHead := latest - ix.cfg.Confirmations

		// 已追平：歇一个轮询间隔再看。/ Caught up: wait one poll interval.
		if from > safeHead {
			if err := ix.sleep(ctx); err != nil {
				return err
			}
			continue
		}

		// 本批区间 [from, to]，最多 BatchSize 个块。/ This batch's range, capped by BatchSize.
		to := from + ix.cfg.BatchSize - 1
		if to > safeHead {
			to = safeHead
		}

		if err := ix.ProcessRange(ctx, from, to); err != nil {
			// 处理失败不推进游标：下一轮重扫同一段（幂等保证不会重复）。
			// On failure don't advance the cursor: re-scan next round (idempotent).
			fmt.Printf("[indexer] 处理 [%d,%d] 失败，将重试 / process failed: %v\n", from, to, err)
			if err := ix.sleep(ctx); err != nil {
				return err
			}
			continue
		}

		// 整批成功落库后才推进游标。/ Advance the cursor only after the whole batch is persisted.
		if err := ix.store.SaveCursor(ctx, ix.cfg.CursorID, to); err != nil {
			return err
		}
		fmt.Printf("[indexer] 已处理 / processed [%d, %d]（链头 %d）\n", from, to, latest)
		from = to + 1

		// 追到安全头了就歇一下，否则立刻扫下一批。/ Pause if we hit the safe head.
		if to == safeHead {
			if err := ix.sleep(ctx); err != nil {
				return err
			}
		}
	}
}

// startBlock 决定从哪个块开始：有游标就续，否则用配置的部署块。
// startBlock decides where to begin: resume from cursor, else the configured deploy block.
func (ix *Indexer) startBlock(ctx context.Context) (uint64, error) {
	last, found, err := ix.store.GetCursor(ctx, ix.cfg.CursorID)
	if err != nil {
		return 0, err
	}
	if found {
		return last + 1, nil
	}
	return ix.cfg.StartBlock, nil
}

// sleep 等待一个轮询间隔，期间若 ctx 取消则立即返回。
// sleep waits one poll interval, returning early if ctx is cancelled.
func (ix *Indexer) sleep(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(ix.cfg.PollInterval):
		return nil
	}
}

// ProcessRange 拉取 [from,to] 的日志并按顺序逐条处理。导出是为了一次性的
// 补录脚本（cmd/backfill-activity）能复用同一套解析/落库逻辑。
// ProcessRange fetches logs in [from,to] and handles them in order. Exported
// so one-off backfill tooling (cmd/backfill-activity) can reuse the same
// parse/persist logic without duplicating it.
func (ix *Indexer) ProcessRange(ctx context.Context, from, to uint64) error {
	logs, err := ix.chain.FilterLogs(ctx, from, to)
	if err != nil {
		return err
	}
	// 节点返回的日志已按 (区块, 日志序号) 升序，顺序处理即符合链上时序。
	// Node returns logs ordered by (block, logIndex); processing in order preserves on-chain time.
	for i := range logs {
		if err := ix.handleLog(ctx, &logs[i]); err != nil {
			return fmt.Errorf("区块 %d 交易 %s: %w", logs[i].BlockNumber, logs[i].TxHash.Hex(), err)
		}
	}
	return nil
}

// ======================= 下半：事件分发与处理 / dispatch & handlers =======================

// handleLog 按 topics[0]（事件签名哈希）把一条原始日志分发给对应处理函数。
// handleLog dispatches one raw log to its handler by topics[0] (the event signature hash).
func (ix *Indexer) handleLog(ctx context.Context, lg *types.Log) error {
	if len(lg.Topics) == 0 {
		return nil // 匿名事件，忽略 / anonymous event, skip
	}
	switch lg.Topics[0] {
	case ix.escrowABI.Events["CampaignCreated"].ID:
		return ix.onCampaignCreated(ctx, lg)
	case ix.escrowABI.Events["DonationReceived"].ID:
		return ix.onDonationReceived(ctx, lg)
	case ix.escrowABI.Events["MilestoneReleased"].ID:
		return ix.onMilestoneReleased(ctx, lg)
	case ix.escrowABI.Events["ReceiptSubmitted"].ID:
		return ix.onReceiptSubmitted(ctx, lg)
	case ix.escrowABI.Events["CampaignCompleted"].ID:
		return ix.onCampaignCompleted(ctx, lg)
	case ix.registryABI.Events["CharityVerified"].ID:
		return ix.onCharity(ctx, lg, true)
	case ix.registryABI.Events["CharityRevoked"].ID:
		return ix.onCharity(ctx, lg, false)
	default:
		return nil // 其它事件（如 OwnershipTransferred）不关心 / ignore the rest
	}
}

// onCampaignCreated：事件只给了 goal/数量，里程碑金额与 metadata 用 view 补齐后一并落库。
// The event only carries goal/count; milestone amounts and metadata are backfilled via view calls.
func (ix *Indexer) onCampaignCreated(ctx context.Context, lg *types.Log) error {
	ev, err := ix.chain.Escrow.ParseCampaignCreated(*lg)
	if err != nil {
		return err
	}
	opts := &bind.CallOpts{Context: ctx}
	msRaw, err := ix.chain.Escrow.GetMilestones(opts, ev.CampaignId) // 补里程碑金额
	if err != nil {
		return err
	}
	camp, err := ix.chain.Escrow.GetCampaign(opts, ev.CampaignId) // 补 metadataHash
	if err != nil {
		return err
	}

	c := store.Campaign{
		ID:             ev.CampaignId.Uint64(),
		Charity:        lowerAddr(ev.Charity),
		Goal:           ev.Goal,
		MilestoneCount: uint32(ev.MilestoneCount.Uint64()),
		MetadataHash:   hashHex(camp.MetadataHash),
		CreatedBlock:   lg.BlockNumber,
		CreatedTx:      lg.TxHash.Hex(),
	}
	ms := make([]store.Milestone, len(msRaw))
	for i, m := range msRaw {
		ms[i] = store.Milestone{CampaignID: c.ID, Idx: uint32(i), Amount: m.Amount}
	}
	if err := ix.store.InsertCampaign(ctx, c, ms, true); err != nil {
		return err
	}
	return ix.store.InsertActivity(ctx, store.ActivityEvent{
		CampaignID: c.ID, EventType: "CampaignCreated",
		BlockNumber: lg.BlockNumber, TxHash: lg.TxHash.Hex(), LogIndex: uint32(lg.Index),
	})
}

func (ix *Indexer) onDonationReceived(ctx context.Context, lg *types.Log) error {
	ev, err := ix.chain.Escrow.ParseDonationReceived(*lg)
	if err != nil {
		return err
	}
	d := store.Donation{
		CampaignID:  ev.CampaignId.Uint64(),
		Donor:       lowerAddr(ev.Donor),
		Amount:      ev.Amount,
		BlockNumber: lg.BlockNumber,
		TxHash:      lg.TxHash.Hex(),
		LogIndex:    uint32(lg.Index),
	}
	if err := ix.store.InsertDonation(ctx, d); err != nil {
		return err
	}
	return ix.store.InsertActivity(ctx, store.ActivityEvent{
		CampaignID: d.CampaignID, EventType: "DonationReceived", Amount: d.Amount,
		BlockNumber: d.BlockNumber, TxHash: d.TxHash, LogIndex: d.LogIndex,
	})
}

func (ix *Indexer) onMilestoneReleased(ctx context.Context, lg *types.Log) error {
	ev, err := ix.chain.Escrow.ParseMilestoneReleased(*lg)
	if err != nil {
		return err
	}
	idx := uint32(ev.MilestoneIndex.Uint64())
	if err := ix.store.MarkMilestoneReleased(ctx, ev.CampaignId.Uint64(), idx); err != nil {
		return err
	}
	return ix.store.InsertActivity(ctx, store.ActivityEvent{
		CampaignID: ev.CampaignId.Uint64(), EventType: "MilestoneReleased",
		Amount: ev.Amount, MilestoneIdx: &idx,
		BlockNumber: lg.BlockNumber, TxHash: lg.TxHash.Hex(), LogIndex: uint32(lg.Index),
	})
}

func (ix *Indexer) onReceiptSubmitted(ctx context.Context, lg *types.Log) error {
	ev, err := ix.chain.Escrow.ParseReceiptSubmitted(*lg)
	if err != nil {
		return err
	}
	idx := uint32(ev.MilestoneIndex.Uint64())
	if err := ix.store.MarkMilestoneProven(ctx, ev.CampaignId.Uint64(), idx, hashHex(ev.ReceiptHash)); err != nil {
		return err
	}
	return ix.store.InsertActivity(ctx, store.ActivityEvent{
		CampaignID: ev.CampaignId.Uint64(), EventType: "ReceiptSubmitted", MilestoneIdx: &idx,
		BlockNumber: lg.BlockNumber, TxHash: lg.TxHash.Hex(), LogIndex: uint32(lg.Index),
	})
}

func (ix *Indexer) onCampaignCompleted(ctx context.Context, lg *types.Log) error {
	ev, err := ix.chain.Escrow.ParseCampaignCompleted(*lg)
	if err != nil {
		return err
	}
	if err := ix.store.MarkCampaignCompleted(ctx, ev.CampaignId.Uint64()); err != nil {
		return err
	}
	return ix.store.InsertActivity(ctx, store.ActivityEvent{
		CampaignID: ev.CampaignId.Uint64(), EventType: "CampaignCompleted",
		BlockNumber: lg.BlockNumber, TxHash: lg.TxHash.Hex(), LogIndex: uint32(lg.Index),
	})
}

// onCharity 处理认证/撤销两个事件（结构相同，只是 verified 不同）。
// onCharity handles both verify/revoke events (same shape, differing verified flag).
func (ix *Indexer) onCharity(ctx context.Context, lg *types.Log, verified bool) error {
	var charity common.Address
	if verified {
		ev, err := ix.chain.Registry.ParseCharityVerified(*lg)
		if err != nil {
			return err
		}
		charity = ev.Charity
	} else {
		ev, err := ix.chain.Registry.ParseCharityRevoked(*lg)
		if err != nil {
			return err
		}
		charity = ev.Charity
	}
	return ix.store.UpsertCharity(ctx, lowerAddr(charity), verified, lg.BlockNumber)
}

// --- 小工具 / helpers ---

// lowerAddr 把地址规范成小写 0x 字符串（与库里一致）。
// lowerAddr normalizes an address to a lowercase 0x string (matching the DB).
func lowerAddr(a common.Address) string { return strings.ToLower(a.Hex()) }

// hashHex 把 bytes32 转成 0x + 64 位十六进制。/ hashHex renders a bytes32 as 0x + 64 hex chars.
func hashHex(b [32]byte) string { return common.BytesToHash(b[:]).Hex() }
