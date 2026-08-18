// Command api serves the backend HTTP API.
//
// On boot it, in order:
//  1. reads DATABASE_URL from the environment,
//  2. applies every pending migration from the embedded migrations/ directory,
//  3. listens on $PORT (falling back to $APP_PORT, then 8080).
//
// /healthz reports 200 only after migrations succeeded and a SELECT 1 against
// the database works, so a healthy check means the app can actually serve.
//
// Routes are registered WITHOUT the /api prefix: the edge proxy routes this
// service with handle_path /api/*, which strips the prefix before forwarding.
// The browser requests /api/v1/content and this server receives /v1/content.
package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := migrate(ctx, db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz(db))
	// No method restriction: the error catalog (services.md §2.3) has no
	// 405/VALIDATION_FAILED code, and the endpoint is read-only with no body,
	// so every method resolves to the same public read.
	mux.HandleFunc("/v1/content", handleContent(db))

	addr := ":" + listenPort()
	log.Printf("listening on %s", addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           withRequestID(withLogging(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func listenPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	if p := os.Getenv("APP_PORT"); p != "" {
		return p
	}
	return "8080"
}

// ── Handlers ────────────────────────────────────────────────────────────────

func handleContent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var value string
		err := db.QueryRowContext(ctx,
			`SELECT value FROM contents LIMIT 1`).Scan(&value)
		if err != nil {
			status, code, message := classifyContentErr(err)
			writeError(w, r, status, code, message)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"value": value})
	}
}

func handleHealthz(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var one int
		if err := db.QueryRowContext(ctx, `SELECT 1`).Scan(&one); err != nil {
			writeError(w, r, http.StatusServiceUnavailable, "UNAVAILABLE",
				"Database is unavailable")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// classifyContentErr maps a read failure to the closed error catalog
// (services.md §2.3): NOT_FOUND (no row), UNAVAILABLE (DB unreachable or
// query timeout), INTERNAL (anything else, e.g. a scan error).
func classifyContentErr(err error) (status int, code, message string) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, "NOT_FOUND", "The content row does not exist"
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		return http.StatusServiceUnavailable, "UNAVAILABLE", "Database is unavailable"
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return http.StatusServiceUnavailable, "UNAVAILABLE", "Database is unavailable"
	}
	var connectErr *pgconn.ConnectError
	if errors.As(err, &connectErr) {
		return http.StatusServiceUnavailable, "UNAVAILABLE", "Database is unavailable"
	}

	return http.StatusInternalServerError, "INTERNAL", "Internal error"
}

// ── HTTP plumbing ───────────────────────────────────────────────────────────

type ctxKey int

const requestIDKey ctxKey = iota

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	Details   []string `json:"details"`
	RequestID string   `json:"request_id"`
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:      code,
			Message:   message,
			Details:   []string{},
			RequestID: requestIDFrom(r.Context()),
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

func requestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("request_id=%s method=%s path=%s status=%d duration_ms=%d",
			requestIDFrom(r.Context()), r.Method, r.URL.Path,
			rec.status, time.Since(start).Milliseconds())
	})
}

// ── Migrations ──────────────────────────────────────────────────────────────

// migrate applies embedded *.up.sql files in filename order, recording each in
// a schema_migrations table so re-running is a no-op. Down files are not
// applied automatically.
func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		name       TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := Files.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	var ups []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			ups = append(ups, e.Name())
		}
	}
	sort.Strings(ups)

	for _, name := range ups {
		var done bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE name = $1)`, name,
		).Scan(&done); err != nil {
			return fmt.Errorf("check %s: %w", name, err)
		}
		if done {
			continue
		}

		body, err := Files.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations (name) VALUES ($1)`, name,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit %s: %w", name, err)
		}
		log.Printf("applied migration %s", name)
	}

	return nil
}
