package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"personal-financial-assistant/backend/internal/auth"
	"personal-financial-assistant/backend/internal/db"
	"personal-financial-assistant/backend/internal/httpapi"
	"personal-financial-assistant/backend/internal/users"
)

func main() {
	addr := env("ADDR", ":8080")

	databaseURL := env("DATABASE_URL", "postgres://ledgerlite:ledgerlite@localhost:5432/ledgerlite?sslmode=disable")
	jwtSecret := env("JWT_SECRET", "dev_secret_change_me")

	ctx := context.Background()

	dbConn, err := db.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	jwtMgr, err := auth.NewManager(jwtSecret, 24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	usersRepo := users.NewRepo(dbConn.Pool)

	handler := httpapi.NewRouter(httpapi.Deps{
		UsersRepo: usersRepo,
		JWT:       jwtMgr,
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
