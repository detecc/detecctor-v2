FROM golang:1.17 as base
WORKDIR /detecctor-v2/notification-service/src
COPY . .

WORKDIR /detecctor-v2/notification-service/src/cmd/notifications

FROM base as dev
ENTRYPOINT ["go","run","."]

FROM base as build
RUN go build main.go -o ../../../notification-service

FROM base as test
RUN go test -v

FROM alpine as app
COPY --from=build /detecctor-v2/notification-service/notification-service /detecctor-v2/notification-service
RUN cp /detecctor-v2/notification-service/notification-service /usr/bin/notification-service

ENTRYPOINT ["notification-service"]