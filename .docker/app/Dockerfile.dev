FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN go build ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates libc6-compat ffmpeg python3 py3-pip

# Install Instaloader for Instagram support using virtual environment
RUN python3 -m venv /opt/instaloader-env && \
    /opt/instaloader-env/bin/pip install instaloader && \
    ln -s /opt/instaloader-env/bin/instaloader /usr/local/bin/instaloader

WORKDIR /app
COPY --from=builder /app .

CMD ["./main"]