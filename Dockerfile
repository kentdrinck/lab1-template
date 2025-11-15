FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/server/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /app
COPY . .
COPY --from=builder /app/bin bin

CMD ["./bin/main"]