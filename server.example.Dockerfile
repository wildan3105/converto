FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /engine ./main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /engine /engine

RUN chmod +x /engine

EXPOSE 3000

CMD ["/engine", "server"]
