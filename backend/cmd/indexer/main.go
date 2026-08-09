// Command indexer 是事件索引器的入口：装配各层并运行主循环，支持优雅退出。
// Command indexer is the event indexer's entrypoint: it wires the layers,
// runs the main loop, and shuts down gracefully on Ctrl+C / SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zoujunhui1/trust-chain/backend/internal/chain"
	"github.com/zoujunhui1/trust-chain/backend/internal/config"
	"github.com/zoujunhui1/trust-chain/backend/internal/indexer"
	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

func main() {
	// context.Canceled 是正常退出（收到信号），不当作错误。
	// context.Canceled is a normal (signal-triggered) exit, not an error.
	if err := run(); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	// 收到 Ctrl+C / SIGTERM 时取消 ctx，各层据此停止。
	// Cancel ctx on Ctrl+C / SIGTERM; every layer stops off this signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	fmt.Println("[main] 索引器启动，Ctrl+C 退出 / indexer started, Ctrl+C to stop")
	if err := ix.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Println("[main] 已优雅退出 / graceful shutdown")
			return nil
		}
		return err
	}
	return nil
}
