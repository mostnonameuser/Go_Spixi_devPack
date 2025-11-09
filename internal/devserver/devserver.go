package devserver

import (
	"log"
	"net/http"
)

type HTTPServer struct {
	addr string
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{addr: addr}
}

func (s *HTTPServer) Start() {
	fs := http.FileServer(http.Dir("../../web/app"))
	http.Handle("/", fs)

	log.Printf("Serving frontend at http://localhost%s", s.addr)
	if err := http.ListenAndServe(s.addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}