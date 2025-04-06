FROM golang:1.24-alpine3.21 AS base
RUN go install github.com/air-verse/air@latest
RUN mkdir -p /app/tmp/air
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
WORKDIR /app
COPY . .
CMD ["air", "-c", ".air.toml"]