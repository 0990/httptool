package main

import (
	"flag"
	"github.com/0990/httptool"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var qrcodeDir = flag.String("qrcode_dir", "qrcode_dir", "qrcode pic dir")
var listen = flag.String("listen", "0.0.0.0:80", "http server listen addr")
var advertisedAddr = flag.String("advertised_addr", "127.0.0.1:80", "advertisedAddr")

func main() {
	flag.Parse()
	parseEnv()

	qrcode, err := httptool.NewQRCode(*qrcodeDir)
	if err != nil {
		logrus.WithError(err).Error("NewQRCode")
		return
	}
	s := httptool.NewServer(*listen, *advertisedAddr, qrcode)
	s.Run()

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		gracefullShutDown(s, quit, done)
	}()

	<-done
}

func parseEnv() {
	if dir := os.Getenv("QRCode_Dir"); len(dir) > 0 {
		qrcodeDir = &dir
	}

	if lis := os.Getenv("Listen"); len(lis) > 0 {
		listen = &lis
	}

	if aAddr := os.Getenv("Advertised_Addr"); len(aAddr) > 0 {
		advertisedAddr = &aAddr
	}
}

func gracefullShutDown(server *httptool.Server, quit <-chan os.Signal, done chan<- bool) {
	sig := <-quit
	logrus.Infof("Server is shutting down exiting... signal:(%v)", sig)
	err := server.GraceShutDown()
	if err != nil {
		logrus.Fatalf("Could not gracefully shutdown the serfver:%v\n", err)
	}

	logrus.Info("Server shut down")
	close(done)
}
