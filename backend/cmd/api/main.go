// Command api 是读侧 HTTP 服务的入口：连库、起服务器、优雅退出。
// Command api is the read-side HTTP server entrypoint: connect the DB,
// start the server, and shut down gracefully on Ctrl+C / SIGTERM.
//
// 它和 cmd/indexer 是两个独立进程 / two separate processes:
//   indexer 只写库（链→MySQL），api 只读库（MySQL→前端），互不阻塞。
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zoujunhui1/trust-chain/backend/internal/api"
	"github.com/zoujunhui1/trust-chain/backend/internal/config"
	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	// 收到 Ctrl+C / SIGTERM 时取消 ctx，用于触发优雅关闭。
	// Cancel ctx on Ctrl+C / SIGTERM to trigger graceful shutdown.
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

	// http.Server 显式配好超时，避免慢连接耗尽资源（默认零值是没有超时的）。
	// Configure timeouts explicitly; the zero-value server has none, which is unsafe.
	srv := &http.Server{
		Addr:              cfg.APIAddr,
		Handler:           api.New(st).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 在独立 goroutine 里跑 ListenAndServe，主 goroutine 等退出信号。
	// Run ListenAndServe in a goroutine; the main goroutine waits for the signal.
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("[api] 监听 / listening on %s\n", cfg.APIAddr)
		// 正常关闭时 ListenAndServe 返回 ErrServerClosed，不算错误。
		// On graceful close it returns ErrServerClosed, which isn't an error.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// 要么服务器自己出错，要么收到退出信号。/ Either the server errors, or we get a signal.
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		fmt.Println("\n[api] 收到退出信号，优雅关闭中 / shutting down gracefully")
	}

	// 给正在处理的请求最多 10 秒收尾，再强制关闭。
	// Give in-flight requests up to 10s to finish before forcing close.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("api: 关闭超时 / shutdown: %w", err)
	}
	fmt.Println("[api] 已退出 / stopped")
	return nil
}
