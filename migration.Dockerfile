FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

COPY migrations /migrations

CMD ["sh", "-c", "sleep 2 && goose -dir \"${MIGRATION_DIR}\" postgres \"${MIGRATION_DSN}\" up"]

