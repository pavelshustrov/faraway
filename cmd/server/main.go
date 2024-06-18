package main

import (
	"context"
	"faraway/internal/server"
	"faraway/internal/store/quotes"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

var config server.Config

func init() {
	flag.StringVar(&config.Port, "PORT", "8080", "Server port")
	flag.StringVar(&config.DDosProtectionStrategy, "DDOS_PROTECTION", "OFF", "Enable or disable DDoS protection (default: OFF)")
	flag.DurationVar(&config.WriteTimeout, "WRITE_TIMEOUT", 0, "Write timeout connection")
	flag.DurationVar(&config.ReadTimeout, "READ_TIMEOUT", 0, "Read timeout connection")

	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	quoteStorage := quotes.NewStore()

	s := server.New(config, quoteStorage)

	if err := s.Serve(ctx); err != nil {
		log.Fatal(err)
	}
}
