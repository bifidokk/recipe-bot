FROM golang:alpine AS builder

RUN apk add --no-cache git

# Set GOROOT explicitly
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENV PATH=$PATH:/go/bin:/root/go/bin

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY migrations/ /migrations/

WORKDIR /

ENTRYPOINT ["goose"]