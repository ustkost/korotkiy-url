FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.24.4
RUN adduser -D app
WORKDIR /app
COPY --from=builder --chown=app:app /app/server .
USER app
EXPOSE 8080
CMD ["./server"]
