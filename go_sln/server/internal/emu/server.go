package emu

import (
	"fmt"
	"log"
	"net"
	"sync"

	"sln/internal/config"
)

type Server struct {
	cfg    *config.Config
	logger *log.Logger
	ln     net.Listener
	wg     sync.WaitGroup
	close  chan struct{}
	closed bool
	mu     sync.Mutex
}

func NewServer(cfg *config.Config, logger *log.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		close:  make(chan struct{}),
	}
}

// запуск TCP-слушатля и прием подключения
// блок до ошибки или Stop()
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.ln = ln
	s.logger.Printf("listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.close:
				// завершение ?
				return nil
			default:
				s.logger.Printf("accept error: %v", err)
				continue
			}
		}
		s.logger.Printf("accepted connection from %s", conn.RemoteAddr())
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			handleConnection(c, s.cfg, s.logger)
		}(conn)
	}
}

// стоп сервера, закрытие прослушивателя и вызов горутин
func (s *Server) Stop() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	s.mu.Unlock()

	close(s.close)
	if s.ln != nil {
		_ = s.ln.Close()
	}
	s.logger.Printf("closing server, waiting for handlers...")
	s.wg.Wait()
	s.logger.Printf("server stopped")
}
