FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg
COPY config/ ./config
COPY docs/ ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/


FROM alpine:3.20

WORKDIR /app
LABEL maintainer="pomapinigin@gmail.com"

COPY --from=builder /app/main ./main
COPY migrations/ ./migrations

EXPOSE 8080
CMD ["./main"]