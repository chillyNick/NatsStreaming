package http_server

import (
	"fmt"
	"net/http"

	"github.com/chillyNick/NatsStreaming/internal/cache"
	"github.com/chillyNick/NatsStreaming/internal/config"
)

type HttpServer struct {
	c      cache.Cache
	server *http.Server
}

func New(cfg config.Rest, c cache.Cache) HttpServer {
	httpServer := HttpServer{c: c}

	mux := http.DefaultServeMux
	mux.HandleFunc("/", httpServer.GetOrder)

	httpServer.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Handler: mux,
	}

	return httpServer
}

func (h HttpServer) Start() error {
	return h.server.ListenAndServe()
}

func (h HttpServer) Close() error {
	return h.server.Close()
}

func (h HttpServer) GetOrder(rw http.ResponseWriter, r *http.Request) {
	info, ok := h.c.Get(r.URL.Query().Get("orderUid"))
	if !ok {
		rw.Write([]byte("Order with such id not found"))

		return
	}

	rw.Write(info)
}
