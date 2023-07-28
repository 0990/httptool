package httptool

import (
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
	"time"
)

type QRCode struct {
	dir string
}

func NewQRCode(dir string) (*QRCode, error) {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &QRCode{dir: dir}, nil
}

func (q *QRCode) Create(content string) (string, error) {
	filename := "qrcode" + time.Now().Format("20060102150405") + ".png"
	fullFileName := q.dir + "/" + filename
	err := qrcode.WriteFile(content, qrcode.Medium, 256, fullFileName)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (q *QRCode) FileServer() http.Handler {
	return http.FileServer(http.Dir(q.dir))
}
