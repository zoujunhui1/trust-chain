// Package chain 封装与以太坊节点（Sepolia）的连接。
// Package chain wraps the connection to an Ethereum node (Sepolia).
//
// 它对上层只暴露索引器真正需要的能力：查最新块高、按区块范围拉原始日志、
// 以及持有 abigen 生成的合约绑定（用于解析事件和调用 view）。
// It exposes only what the indexer needs: latest block height, raw logs by
// block range, and the abigen bindings (to parse events and call view methods).
package chain

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/zoujunhui1/trust-chain/backend/internal/config"
	"github.com/zoujunhui1/trust-chain/backend/internal/contracts"
)

// Client 是对 ethclient 的轻封装，附带两个合约绑定。
// Client is a thin wrapper over ethclient plus the two contract bindings.
type Client struct {
	eth *ethclient.Client

	// Registry / Escrow：abigen 绑定，供上层 Parse 事件与调用 view。
	// The abigen bindings; upper layers use them to parse events and call views.
	Registry *contracts.CharityRegistry
	Escrow   *contracts.TrustChainEscrow

	// addresses：只订阅这两个合约产生的日志。
	// addresses: only subscribe to logs emitted by these two contracts.
	addresses []common.Address
}

// Dial 连接节点、校验连通性，并初始化合约绑定。
// Dial connects to the node, checks connectivity, and initializes the bindings.
func Dial(ctx context.Context, cfg *config.Config) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("chain: 连接 RPC 失败 / dial RPC: %w", err)
	}

	// 立刻探一次链 ID，确认 RPC 真的可用（fail fast）。
	// Probe the chain ID right away to confirm the RPC works (fail fast).
	chainID, err := eth.ChainID(ctx)
	if err != nil {
		eth.Close()
		return nil, fmt.Errorf("chain: 读取 chainID 失败 / get chainID: %w", err)
	}

	registry, err := contracts.NewCharityRegistry(cfg.RegistryAddress, eth)
	if err != nil {
		eth.Close()
		return nil, fmt.Errorf("chain: 绑定 registry 失败 / bind registry: %w", err)
	}
	escrow, err := contracts.NewTrustChainEscrow(cfg.EscrowAddress, eth)
	if err != nil {
		eth.Close()
		return nil, fmt.Errorf("chain: 绑定 escrow 失败 / bind escrow: %w", err)
	}

	fmt.Printf("[chain] 已连接 / connected  chainID=%s  registry=%s  escrow=%s\n",
		chainID, cfg.RegistryAddress.Hex(), cfg.EscrowAddress.Hex())

	return &Client{
		eth:       eth,
		Registry:  registry,
		Escrow:    escrow,
		addresses: []common.Address{cfg.RegistryAddress, cfg.EscrowAddress},
	}, nil
}

// LatestBlock 返回当前链头区块高度。/ LatestBlock returns the current head block height.
func (c *Client) LatestBlock(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

// FilterLogs 拉取 [from, to]（闭区间）内、来自我们两个合约的所有原始日志。
// FilterLogs fetches all raw logs from our two contracts within [from, to] (inclusive).
// 上层再按 topic 分发、用绑定的 Parse* 解析成强类型事件。
// The caller then dispatches by topic and parses via the bindings' Parse* methods.
func (c *Client) FilterLogs(ctx context.Context, from, to uint64) ([]types.Log, error) {
	q := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: c.addresses,
	}
	logs, err := c.eth.FilterLogs(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("chain: FilterLogs [%d,%d] 失败: %w", from, to, err)
	}
	return logs, nil
}

// Close 关闭底层连接。/ Close releases the underlying connection.
func (c *Client) Close() { c.eth.Close() }
