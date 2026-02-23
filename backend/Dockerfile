# 1. Stage: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Wir kopieren alles aus dem Kontext (der durch Compose definiert wird)
COPY . .

# Falls go.mod im Hauptordner liegt, aber das Dockerfile im Unterordner, 
# müssen wir sicherstellen, dass die Pfade stimmen. 
# Wir laden die Module und bauen die App.
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./backend/cmd/server/main.go

# 2. Stage: Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 1323
CMD ["./main"]
