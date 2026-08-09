// Package config 集中读取并校验索引器的运行配置。
// Package config loads and validates the indexer's runtime configuration.
//
// 所有环境相关的输入（RPC、数据库、合约地址、起始区块等）只在这里读一次，
// 其它模块拿到的是一个已校验、强类型的 Config，不再直接碰 os.Getenv。
// All env-dependent inputs are read once here; other packages receive a
// validated, strongly-typed Config instead of touching os.Getenv directly.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

// Config 是索引器的完整运行配置。/ Config is the indexer's full runtime configuration.
type Config struct {
	RPCURL          string         // Sepolia JSON-RPC 端点（Alchemy）
	MySQLDSN        string         // MySQL 连接串 / DSN
	RegistryAddress common.Address // CharityRegistry 合约地址
	EscrowAddress   common.Address // TrustChainEscrow 合约地址
	StartBlock      uint64         // 首次回填的起始区块（合约部署块）
	CursorID        string         // 游标主键，标识这条链，例如 "sepolia"
	PollInterval    time.Duration  // 追到链头后，隔多久再拉一次新块
	Confirmations   uint64         // 落后链头多少块再处理，规避链重组
	BatchSize       uint64         // 每次 FilterLogs 最多扫多少个区块
	APIAddr         string         // REST API 监听地址，如 ":8080"（仅 api 二进制用）
}

// Load 从进程环境（可选先加载 .env 文件）读取并校验配置。
// Load reads and validates configuration from the process environment,
// optionally loading a .env file first.
func Load() (*Config, error) {
	// .env 只是「本地便利」：不存在也没关系，真实环境用系统环境变量。
	// The .env file is a local convenience; absence is fine (prod uses real env vars).
	_ = godotenv.Load()

	cfg := &Config{
		RPCURL:        mustEnv("RPC_URL"),
		MySQLDSN:      mustEnv("MYSQL_DSN"),
		CursorID:      envOr("CURSOR_ID", "sepolia"),
		PollInterval:  time.Duration(envUintOr("POLL_INTERVAL_SECONDS", 12)) * time.Second,
		Confirmations: envUintOr("CONFIRMATIONS", 5),
		BatchSize:     envUintOr("BATCH_SIZE", 2000),
		APIAddr:       envOr("API_ADDR", ":8080"),
	}

	// 合约地址：必须是合法 0x 地址，否则直接报错（fail fast）。
	// Contract addresses must be valid 0x addresses, otherwise fail fast.
	var err error
	if cfg.RegistryAddress, err = parseAddress("REGISTRY_ADDRESS"); err != nil {
		return nil, err
	}
	if cfg.EscrowAddress, err = parseAddress("ESCROW_ADDRESS"); err != nil {
		return nil, err
	}

	if cfg.StartBlock, err = parseUint("START_BLOCK"); err != nil {
		return nil, err
	}

	if cfg.BatchSize == 0 {
		return nil, fmt.Errorf("config: BATCH_SIZE 必须大于 0 / must be > 0")
	}
	return cfg, nil
}

// --- 下面是读取环境变量的小工具 / helpers ---

// mustEnv 读取必填变量，缺失时 panic（配置错误应尽早、明确地暴露）。
// mustEnv reads a required var, panicking if missing (config errors should surface early).
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: 缺少必填环境变量 %s / missing required env var %s", key, key))
	}
	return v
}

// envOr 读取可选变量，缺失时返回默认值。/ envOr reads an optional var with a fallback.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// envUintOr 读取可选的无符号整数，缺失或非法时返回默认值。
// envUintOr reads an optional unsigned int, falling back on missing/invalid input.
func envUintOr(key string, fallback uint64) uint64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

// parseUint 读取必填的无符号整数。/ parseUint reads a required unsigned int.
func parseUint(key string) (uint64, error) {
	v := mustEnv(key)
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("config: %s 不是合法整数 / not a valid integer: %q", key, v)
	}
	return n, nil
}

// parseAddress 读取并校验一个必填的 0x 以太坊地址。
// parseAddress reads and validates a required 0x Ethereum address.
func parseAddress(key string) (common.Address, error) {
	v := mustEnv(key)
	if !common.IsHexAddress(v) {
		return common.Address{}, fmt.Errorf("config: %s 不是合法地址 / not a valid address: %q", key, v)
	}
	return common.HexToAddress(v), nil
}
