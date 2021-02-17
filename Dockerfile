FROM golang:1.15-alpine as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

WORKDIR /build

COPY go.mod .
# COPY go.sum .

RUN go mod download

COPY . .

# RUN go build cmd/server/main.go -o /build/waveformServer .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM yunfei1982/ub-waveform
WORKDIR /app
COPY --from=builder /build/main ./waveformServer

# EXPOSE 9527
ENTRYPOINT ["/app/waveformServer"]