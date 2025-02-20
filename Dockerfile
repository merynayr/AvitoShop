FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/bin/app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/app .
COPY .env .

CMD ["./app", "-config-path=.env"]
