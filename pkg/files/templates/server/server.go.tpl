/*Generated code do not modify it*/
package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Server struct {
	server   *http.Server
	serverWG *sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		server:   &http.Server{},
		serverWG: &sync.WaitGroup{},
	}
}

func (s *Server) SetHandler(router http.Handler) *Server {
	s.server.Handler = router
	return s
}

func (s *Server) SetPort(port int64) *Server {
	s.server.Addr = fmt.Sprintf(":%d", port)
	return s
}

func (s *Server) Start() {
	go func() {
		s.serverWG.Add(1)
		defer s.serverWG.Done()

		log.Infof("Starting server at address %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorf("Server failed: %w", err)
		}
	}()
}

func (s *Server) Stop(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}

	log.Infof("stopping the http server")
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Errorf("failed to shutdown the http server: %s", err)
	}

	s.serverWG.Wait()
}
