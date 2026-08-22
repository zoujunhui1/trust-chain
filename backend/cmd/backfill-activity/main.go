// Command backfill-activity 是一次性脚本：给 activity_events 补上表建好之前
// 就已发生的历史事件（主循环的游标已经跑过那些区块，不会再重新处理它们）。
// 用完可以删掉；重跑是安全的——写入用 INSERT IGNORE，幂等。
//
// Command backfill-activity is a one-off script: it backfills activity_events
// for history that happened before the table existed (the main loop's cursor
// has already passed those blocks and won't reprocess them). Safe to delete
// once run; safe to rerun — writes are idempotent (INSERT IGNORE).
//
// 用法 / usage:
//
//	go run ./cmd/backfill-activity -from 11455513 -to 11455540
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/zoujunhui1/trust-chain/backend/internal/chain"
	"github.com/zoujunhui1/trust-chain/backend/internal/config"
	"github.com/zoujunhui1/trust-chain/backend/internal/indexer"
	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

func main() {
	from := flag.Uint64("from", 0, "起始区块（含）/ start block (inclusive)")
	to := flag.Uint64("to", 0, "结束区块（含）/ end block (inclusive)")
	flag.Parse()

	if *from == 0 || *to == 0 || *to < *from {
		fmt.Fprintln(os.Stderr, "usage: backfill-activity -from <block> -to <block>")
		os.Exit(1)
	}

	if err := run(*from, *to); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run(from, to uint64) error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	st, err := store.Open(cfg.MySQLDSN)
	if err != nil {
		return err
	}
	defer st.Close()

	ch, err := chain.Dial(ctx, cfg)
	if err != nil {
		return err
	}
	defer ch.Close()

	ix, err := indexer.New(cfg, ch, st)
	if err != nil {
		return err
	}

	fmt.Printf("[backfill] 补录活动流水 [%d, %d]\n", from, to)
	if err := ix.ProcessRange(ctx, from, to); err != nil {
		return err
	}
	fmt.Println("[backfill] 完成 / done")
	return nil
}
