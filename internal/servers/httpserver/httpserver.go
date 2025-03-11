package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(ctx context.Context, port int, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
	}
}

func (s *HTTPServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
