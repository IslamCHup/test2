
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app"]