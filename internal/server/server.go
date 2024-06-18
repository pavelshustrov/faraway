package server

import (
	"context"
	"errors"
	"faraway/internal/proof_of_work"
	"faraway/internal/store/quotes"
	"log"
	"net"
	"sync"
)

type Server struct {
	cfg          Config
	handlersPool *sync.Pool
}

func New(cfg Config, storage *quotes.QuoteStore) *Server {
	pool := &sync.Pool{
		New: func() interface{} {
			return &ConnectionHandler{
				readBuf: make([]byte, 1024),
				quotes:  storage,
				pow:     getDDosProtection(cfg.DDosProtectionStrategy),
			}
		},
	}

	return &Server{
		cfg:          cfg,
		handlersPool: pool,
	}
}

func (server *Server) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+server.cfg.Port)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		handler := server.handlersPool.Get().(*ConnectionHandler)

		go func(conn net.Conn, handler *ConnectionHandler) {
			defer func() {
				handler.readBuf = handler.readBuf[:0]
				server.handlersPool.Put(handler)
			}()

			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
				log.Println("closing connection")
				if err := conn.Close(); err != nil {
					log.Println(err)
				}
			}()

			if err := handler.Handle(NewTimeoutConn(conn, server.cfg.ReadTimeout, server.cfg.WriteTimeout)); err != nil {
				if errors.Is(err, proof_of_work.ErrVerificationFailed) {
					log.Println("verification failure from", conn.RemoteAddr())
				}
				return
			}
		}(conn, handler)
	}
}
