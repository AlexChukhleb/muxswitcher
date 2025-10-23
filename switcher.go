package muxswitcher

import (
	"errors"
	"net/http"
	"sync/atomic"
)

var ErrNotUpdated = errors.New("handler not updated")

type Switcher interface {
	http.Handler
	NewHandler(handler http.Handler) error
}

type switcher struct {
	handler         atomic.Value
	handlerNotFound http.Handler
}

func New(handler, notFoundHandler http.Handler) *switcher {
	if notFoundHandler == nil {
		notFoundHandler = http.NotFoundHandler()
	}

	s := &switcher{
		handlerNotFound: notFoundHandler,
	}

	_ = s.NewHandler(handler)

	return s
}

func (s *switcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := s.handler.Load()
	if v == nil {
		s.handlerNotFound.ServeHTTP(w, r)
		return
	}

	h, ok := v.(http.Handler)
	if !ok || h == nil {
		s.handlerNotFound.ServeHTTP(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func (s *switcher) NewHandler(handler http.Handler) error {
	if handler == nil {
		return ErrNotUpdated
	}

	var h http.Handler = handler
	s.handler.Store(h)

	return nil
}
