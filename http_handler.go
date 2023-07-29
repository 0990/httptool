package httptool

import (
	"fmt"
	"net"
	"net/http"
)

func (s *Server) handleQRCode(code *QRCode) http.HandlerFunc {
	type response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		PicUrl  string `json:"picUrl"`
	}

	type request struct {
		Content string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := s.decode(w, r, &req)
		if err != nil {
			s.respond(w, r, response{Code: 3, Message: "invalid param"}, http.StatusOK)
			return
		}

		fileName, err := code.Create(req.Content)
		if err != nil {
			s.respond(w, r, response{Code: 4, Message: err.Error()}, http.StatusOK)
			return
		}

		url := fmt.Sprintf("%s/%s/%s", s.advertisedAddr, "qrcode", fileName)
		s.respond(w, r, response{Code: 0, PicUrl: url}, http.StatusOK)
	}
}

func (s *Server) handleMyip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.respondRaw(w, r, err.Error(), http.StatusOK)
			return
		}

		for m, v := range r.Header {
			fmt.Println(m, v)
		}
		s.respondRaw(w, r, host, http.StatusOK)
	}
}
