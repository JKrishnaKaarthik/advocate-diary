package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JKrishnaKaarthik/advocate-diary/internal/app"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/config"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/httpapi"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/mail"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/session"
	"github.com/JKrishnaKaarthik/advocate-diary/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	objects := storage.New(cfg.S3Endpoint, cfg.S3Region, cfg.S3Bucket, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3PathStyle)
	store := app.NewStore(pool, objects)
	if err := store.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	var sender mail.Sender = mail.LogSender{}
	if cfg.SMTPHost != "" {
		sender = mail.SMTPSender{
			Host: cfg.SMTPHost, Port: cfg.SMTPPort, Username: cfg.SMTPUser, Password: cfg.SMTPPassword, From: cfg.SMTPFrom,
		}
		log.Printf("reminder email via %s", cfg.SMTPHost)
	} else {
		log.Printf("SMTP is not configured; reminder emails are written to the log")
	}

	go func() {
		if err := store.Tick(ctx, sender); err != nil && ctx.Err() == nil {
			log.Printf("reminder tick: %v", err)
		}
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := store.Tick(ctx, sender); err != nil && ctx.Err() == nil {
					log.Printf("reminder tick: %v", err)
				}
			}
		}
	}()

	mux := http.NewServeMux()
	httpapi.Register(mux, store, session.New(cfg.SessionSecret, cfg.CookieSecure))
	if st, err := os.Stat(cfg.WebRoot); err == nil && st.IsDir() {
		mux.Handle("/", spa(cfg.WebRoot))
		log.Printf("serving UI from %s", cfg.WebRoot)
	}

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()
	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func spa(dir string) http.Handler {
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		clean := filepath.Clean("/" + r.URL.Path)
		if strings.Contains(filepath.Base(clean), ".") {
			files.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
	})
}
