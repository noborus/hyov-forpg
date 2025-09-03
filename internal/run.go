package internal

import (
	"bufio"
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func Run(query string, nopager bool) {
	connStr := viper.GetString("db.connection")
	db, err := connectDB(connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if query == "" {
		log.Fatal("query is required. Use --query flag.")
	}

	// Create context that cancels on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received interrupt signal, cancelling query...")
		cancel()
	}()

	rows, err := executeQuery(ctx, db, query)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("Query cancelled by user")
			return
		}
		log.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	if nopager {
		// Output to stdout
		if err := printResults(ctx, os.Stdout, rows, false); err != nil {
			if ctx.Err() != nil {
				log.Println("Output cancelled by user")
				return
			}
			log.Fatalf("failed to print results: %v", err)
		}
		return
	}

	// Use oviewer
	r, w := io.Pipe()
	ov, err := setupOviewer(r)
	if err != nil {
		log.Fatalf("failed to setup oviewer: %v", err)
	}
	configureOviewer(ov)
	defer ov.Close()

	go print(ctx, w, rows)

	if err := ov.Run(); err != nil {
		if ctx.Err() != nil {
			log.Println("Oviewer cancelled by user")
			return
		}
		log.Fatalf("oviewer run failed: %v", err)
	}
}

func print(ctx context.Context, w io.WriteCloser, rows *sql.Rows) {
	defer w.Close()
	bufw := bufio.NewWriter(w)

	// Check if context is already cancelled
	if ctx.Err() != nil {
		return
	}

	if err := printResults(ctx, bufw, rows, true); err != nil {
		if ctx.Err() != nil {
			// Context was cancelled during printing
			return
		}
		log.Printf("print error: %v", err)
		return
	}

	if err := bufw.Flush(); err != nil {
		log.Printf("flush error: %v", err)
	}
}
