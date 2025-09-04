package internal

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

var ErrQueryRequired = errors.New("query is required")

// Run executes the provided SQL query and handles output based on the nopager flag.
func Run(query string, nopager bool) error {
	if query == "" {
		return ErrQueryRequired
	}

	ctx := context.Background()

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return runViewResults(ctx, db, query, nopager)
}

// runViewResults executes the query and manages output display, handling context cancellation and resource cleanup.
func runViewResults(ctx context.Context, db *sql.DB, query string, nopager bool) error {
	ctx, cancel := context.WithCancel(ctx)
	setupSignalHandler(cancel)

	rows, err := runQuery(ctx, db, query)
	if err != nil {
		cancel()
		return err
	}

	defer func() {
		// Ensure resources are cleaned up.
		cancel()
		rows.Close()
	}()

	if nopager {
		return outputToStdout(ctx, rows)
	}

	return setupOviewerAndPrint(ctx, rows)
}

func openDB() (*sql.DB, error) {
	connStr := viper.GetString("db.connection")
	db, err := connectDB(connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func runQuery(ctx context.Context, db *sql.DB, query string) (*sql.Rows, error) {
	return executeQuery(ctx, db, query)
}

func setupSignalHandler(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received interrupt signal, cancelling query...")
		cancel()
	}()
}

func outputToStdout(ctx context.Context, rows *sql.Rows) error {
	return printResults(ctx, os.Stdout, rows, false)
}

func setupOviewerAndPrint(ctx context.Context, rows *sql.Rows) error {
	r, w := io.Pipe()
	ov, err := setupOviewer(r)
	if err != nil {
		return err
	}
	configureOviewer(ov)
	defer ov.Close()

	go outputRowsToWriter(ctx, w, rows)

	return ov.Run()
}

func outputRowsToWriter(ctx context.Context, w io.WriteCloser, rows *sql.Rows) {
	defer w.Close()
	bufw := bufio.NewWriter(w)

	if err := printResults(ctx, bufw, rows, true); err != nil {
		log.Printf("print error: %v", err)
		return
	}

	if err := bufw.Flush(); err != nil {
		log.Printf("flush error: %v", err)
	}
}
