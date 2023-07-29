FROM golang:1.20.5 AS builder
COPY . /httptool
WORKDIR /httptool

RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN CGO_ENABLED=0 go build -o /bin/httptool ./cmd/main.go

FROM scratch as httptool
WORKDIR /httptool
COPY --from=builder /bin/httptool .
COPY --from=builder /httptool/html ./html
CMD ["./httptool"]
