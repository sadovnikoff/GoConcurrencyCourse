package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/common"
	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/concurrency"
)

const (
	defaultIdleTimeout    = 300 * time.Second
	defaultBufSize        = 4096
	defaultMaxConnections = 100
)

type database interface {
	HandleQuery(request string) (string, error)
}

// Server - TCP server
type Server struct {
	lis       net.Listener
	semaphore *concurrency.Semaphore

	db          database
	logger      *common.Logger
	idleTimeout time.Duration
	bufferSize  int
}

// NewServer - returns *Server
func NewServer(cfg *common.Config, db database, logger *common.Logger) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config is invalid")
	}

	if db == nil {
		return nil, errors.New("DB is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	listener, err := net.Listen("tcp", cfg.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("server: failed to listen: %v", err)
	}

	timeout, err := time.ParseDuration(cfg.Network.IdleTimeout)
	if err != nil {
		timeout = defaultIdleTimeout
	}

	bufSize, err := common.ParseSize(cfg.Network.MaxMsgSize)
	if err != nil {
		bufSize = defaultBufSize
	}

	maxConnections := cfg.Network.MaxConnections
	if maxConnections == 0 {
		maxConnections = defaultMaxConnections
	}

	srv := &Server{
		lis:       listener,
		semaphore: concurrency.NewSemaphore(maxConnections),

		db:          db,
		logger:      logger,
		idleTimeout: timeout,
		bufferSize:  bufSize,
	}

	return srv, nil
}

// Run - serve
func (s *Server) Run() {
	defer func() {
		if err := s.lis.Close(); err != nil {
			s.logger.Error("failed to close listener %s", err.Error())
		}
		s.logger.Debug("successfully closed listener %s", s.lis.Addr().String())
	}()

	fmt.Printf("In-memory key-value DB server is running on %s\n", s.lis.Addr().String())
	for {
		conn, err := s.lis.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			s.logger.Error("failed to accept %s", err.Error())
			continue
		}

		s.semaphore.Acquire()
		go func(conn net.Conn) {
			defer func() {
				s.semaphore.Release()
			}()

			s.handle(conn)
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("failed to close connection %s", err.Error())
		}
	}()

	request := make([]byte, s.bufferSize)
	for {
		if s.idleTimeout != 0 {
			if err := conn.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Error("failed to set deadline %s", err.Error())
				return
			}
		}

		count, err := conn.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Error("failed to read request: %s", err.Error())
			break
		} else if count == s.bufferSize {
			s.logger.Error("too small buffer size")
			break
		}

		sanitizedRequest := strings.TrimSpace(string(request[:count]))
		response, err := s.db.HandleQuery(sanitizedRequest)
		if err != nil {
			s.logger.Debug("failed to handle query: %s", err.Error())
			response = err.Error()
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			s.logger.Error("failed to write response: %s", err.Error())
			break
		}
	}
}
