FROM golang:1.24-alpine3.21 AS base
#ARG UID=1000
#ARG GID=1000
#ARG USERNAME=nseven
# Crée groupe uniquement si GID non pris, puis crée l'utilisateur avec le GID
#RUN if ! grep ":x:${GID}:" /etc/group >/dev/null; then addgroup -g ${GID} ${USERNAME}; fi && \
#    adduser -D -u ${UID} -g ${GID} ${USERNAME}
# add git
RUN apk add --no-cache git
# add air
RUN go install github.com/air-verse/air@latest
# Préparer dossier de travail avec les bons droits
#RUN mkdir -p /app/runtime/go-mod /app/runtime/air && chown -R ${UID}:${GID} /app
RUN mkdir -p /app/runtime/air
#ENV GOMODCACHE=/app/runtime/go-mod
WORKDIR /app
#USER $USERNAME
#COPY --chown=$UID:$GID go.mod go.sum ./
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
WORKDIR /app
#COPY --chown=$UID:$GID . .
COPY . .
CMD ["air", "-c", ".air.toml"]

FROM base AS build
WORKDIR /app
#COPY --chown=$UID:$GID . .
COPY . .
RUN go build -o dist/hestia main.go

FROM alpine:3.21 AS prod
WORKDIR /app
# certificats (pour les appels HTTPS éventuels)
RUN apk add --no-cache ca-certificates
# uniquement le binaire
COPY --from=build /app/dist/hestia .
CMD ["./todo_formation"]
