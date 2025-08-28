package main

import (
	"os"
	"os/signal"
	"sln/internal/config"
	"sln/internal/emu"
	"sln/internal/logging"
	"syscall"
)

func main() {
	cfg := config.Load()

	logger := logging.New(cfg.LogFile)

	logger.Printf("starting ttp20 emulator (host=%s port=%d crc=%s adapter=%d)",
		cfg.Host, cfg.Port, cfg.CRCMode, cfg.AdapterAddr)

	srv := emu.NewServer(cfg, logger)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case sigv := <-sig:
		logger.Printf("received signal %v, shutting down...", sigv)
		srv.Stop()
	case err := <-errCh:
		if err != nil {
			logger.Printf("server stopped with error: %v", err)
		} else {
			logger.Printf("server stopped")
		}
	}

	logger.Println("bye")
}
