FROM golang:1.24-alpine3.21 AS base
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN mkdir -p /app/tmp/air
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
WORKDIR /app
COPY . .
CMD ["air", "-c", ".air.toml"]