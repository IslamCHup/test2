package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type Device struct {
	ID       int64  `json:"id"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

func initDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost user=app dbname=devices sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}

func deviceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// контекст запроса + таймаут
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		go func(ctx context.Context) {
			select {
			case <-time.After(5 * time.Second):
				log.Println("long debug operation finished")
			case <-ctx.Done():
				log.Println("debug operation canceled:", ctx.Err())
			}
		}(ctx)

		var d Device

		err = db.QueryRowContext(
			ctx,
			`SELECT id, hostname, ip FROM devices WHERE id = $1`,
			id,
		).Scan(&d.ID, &d.Hostname, &d.IP)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "device not found", http.StatusNotFound)
				return
			}
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		// логируем, но не ломаем ответ пользователю
		if _, err := db.ExecContext(
			ctx,
			`INSERT INTO audit_log(device_id, ts, action) VALUES ($1, now(), 'view')`,
			d.ID,
		); err != nil {
			log.Printf("audit log error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(d); err != nil {
			log.Printf("encode error: %v", err)
		}
	}
}

func main() {
	ctx := context.Background()

	db, err := initDB(ctx)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/device", deviceHandler(db))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// запуск сервера
	go func() {
		log.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
