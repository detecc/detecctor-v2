FROM golang:1.17-alpine as base

WORKDIR /detecctor-v2/notification-service/src
COPY .. .

RUN go mod download
RUN go mod verify

RUN go build -o ../cmd-service ./cmd/cmd/main.go

FROM base as dev
WORKDIR /detecctor-v2/notification-service/src/cmd/notifications
ENTRYPOINT ["go","run","."]

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=build /detecctor-v2/notification-service/notification-service /usr/bin/detecctor-v2/notification-service
WORKDIR /usr/bin/detecctor-v2/

ENTRYPOINT ["./notification-service"]