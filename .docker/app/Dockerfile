FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN go build ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates libc6-compat ffmpeg

WORKDIR /app
COPY --from=builder /app .

CMD ["./main"]