package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/arieffadhlan/go-fitbyte/internal/app"
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/database"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitDBConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDBConnection(db)

	appServer := app.NewServer(cfg, db)
	go func() {
		if err := appServer.Listen(cfg.App.Port); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	if err := appServer.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exited")
}
