#BUILD GO APP
FROM golang:1.24.0-bookworm AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLE=1 GOOS=linux go build -o /boilerplate-go ./cmd/api/main.go

# SETUP CONTAINER RELEASE
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /boilerplate-go /boilerplate-go

EXPOSE 8080

USER root:root

ENTRYPOINT ["/boilerplate-go"]
