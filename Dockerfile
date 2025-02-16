FROM golang:1.23-alpine AS builder

COPY . /github.com/merynayr/AvitoShop/source/
WORKDIR /github.com/merynayr/AvitoShop/source/

RUN go mod download
RUN go build -o ./bin/app cmd/main.go

FROM alpine:latest


WORKDIR /root/

COPY --from=builder /github.com/merynayr/AvitoShop/source/bin/app .

COPY .env .

CMD ["./app", "-config-path=.env"]