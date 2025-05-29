#BUILD GO APP
FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./

RUN CGO_ENABLE=1 GOOS=linux go build -o /boilerplate-api ./cmd/api/main.go

# BUILD STATIC BINARY
FROM busybox:1.37 as busybox

# SETUP CONTAINER RELEASE
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /boilerplate-api /boilerplate-api

COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/printenv /bin/printenv
COPY --from=busybox /bin/clear /bin/clear
COPY --from=busybox /bin/echo /bin/echo
COPY --from=busybox /bin/grep /bin/grep
COPY --from=busybox /bin/cat /bin/cat
COPY --from=busybox /bin/ps /bin/ps

EXPOSE 8080

ENTRYPOINT ["/boilerplate-api"]
