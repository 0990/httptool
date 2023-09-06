package httptool

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

type Server struct {
	router         *http.ServeMux
	server         *http.Server
	listen         string
	advertisedAddr string //广播地址
}

func NewServer(listen string, advertisedAddr string, code *QRCode) *Server {
	s := &Server{}
	s.listen = listen
	s.advertisedAddr = advertisedAddr
	s.router = http.NewServeMux()
	s.routes(code)
	return s
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listen)
	if err != nil {
		return err
	}
	s.server = &http.Server{Handler: s.router}
	go func() {
		err := s.server.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"ListenAddr": s.listen,
			}).WithError(err).Error("http.ListenAndServer error")
			panic(err)
		}
	}()
	return nil
}

func (s *Server) GraceShutDown() error {
	//给服务端最多30秒的关闭时间，如果服务端还没关好，强制结束整个服务
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) routes(code *QRCode) {
	s.router.Handle("/", http.FileServer(http.Dir("html")))
	s.router.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("file_dir"))))
	s.router.Handle("/qrcode/create", s.handleQRCode(code))
	s.router.Handle("/qrcode/", http.StripPrefix("/qrcode", code.FileServer()))
	s.router.Handle("/myip", s.handleMyip())
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.WithField("data", data).WithError(err).Error("Server respond json encode error")
		}
	}
}

func (s *Server) respondRaw(w http.ResponseWriter, r *http.Request, data string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(data))
}
