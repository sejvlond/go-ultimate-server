package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	lgr          LOGGER
	wg           *sync.WaitGroup
	shutdownChan chan struct{}
	cfg          *ServerConfig
	mux          *http.ServeMux
	CallShutdown func()
}

func NewServer(lgr LOGGER, cfg *ServerConfig, shutdownFunc func(),
	shutdownChan chan struct{}, wg *sync.WaitGroup) (*Server, error) {

	server := &Server{
		lgr:          lgr,
		wg:           wg,
		shutdownChan: shutdownChan,
		cfg:          cfg,
		CallShutdown: shutdownFunc,
		mux:          http.NewServeMux(),
	}
	if err := server.InitHandlers(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) InitHandlers() error {
	homepage, err := NewHomepageHandler(s.lgr.WithField("handler", "Homepage"))
	if err != nil {
		return err
	}
	s.mux.Handle("/", homepage)

	return nil
}

func (s *Server) run() net.Listener {
	s.lgr.Infof("Starting public interface :%v", s.cfg.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Port))
	if err != nil {
		s.lgr.Errorf("Could not create listener %q", err)
		return nil
	}
	s.wg.Add(1)
	go func() {
		s.lgr.Infof("Starting to serve on %v...", s.cfg.Port)
		err = http.Serve(listener, s.mux)
		s.lgr.Infof("Serving stopped with %q", err)
		s.wg.Done()
	}()
	return listener
}

func (s *Server) Run() {
	s.lgr.Infof("starting")
	listener := s.run()

	if listener == nil {
		s.CallShutdown()
	} else {
		<-s.shutdownChan
	}

	if listener != nil {
		err := listener.Close()
		if err != nil {
			s.lgr.Errorf("Error closing listener '%v'", err)
		}
	}
	s.wg.Done()
	s.lgr.Infof("stopped")
}
