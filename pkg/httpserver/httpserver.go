package httpserver

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	_defaultAddr              = ":80"
	_defaultReadTimeout       = 5 * time.Second
	_defaultWriteTimeout      = 5 * time.Second
	_defaultReadHeaderTimeout = 2 * time.Second
	_defaultIdleTimeout       = 60 * time.Second
	_defaultMaxHeaderBytes    = 1 << 20
	_defaultShutdownTimeout   = 3 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:              _defaultAddr,
		Handler:           handler,
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
		ReadHeaderTimeout: _defaultReadHeaderTimeout,
		IdleTimeout:       _defaultIdleTimeout,
		MaxHeaderBytes:    _defaultMaxHeaderBytes,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	log.Printf("Starting HTTP server on %s", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.notify <- err
		} else {
			s.notify <- nil
		}
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
