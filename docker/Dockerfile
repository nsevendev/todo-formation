FROM golang:1.24-alpine3.21 AS base
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN mkdir -p /app/tmp/air
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

FROM base AS dev
WORKDIR /app
COPY . .
CMD ["air", "-c", ".air.toml"]

FROM base AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/todof main.go
CMD ["ls", "-l", "/app/dist/todof"]

FROM alpine:3.21 AS prod
WORKDIR /app
# certificats (pour les appels HTTPS éventuels)
RUN apk add --no-cache ca-certificates
# uniquement le binaire
COPY --from=build /app/dist/todof .
COPY go.mod .
RUN chmod +x /app/todof
CMD ["./todof"]