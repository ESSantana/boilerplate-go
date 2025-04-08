#BUILD GO APP
FROM golang:1.24.0-bookworm AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./

RUN CGO_ENABLE=1 GOOS=linux go build -o /ella-api ./cmd/api/main.go

# SETUP CONTAINER RELEASE
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /ella-api /ella-api

EXPOSE 8080

USER root:root

ENTRYPOINT ["/ella-api"]
