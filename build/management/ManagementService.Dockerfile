FROM golang:1.17-alpine as base

WORKDIR /detecctor-v2/management-service/src
COPY .. .

RUN go mod download
RUN go mod verify

RUN go build -o ../cmd-service ./cmd/cmd/main.go

FROM base as dev
WORKDIR /detecctor-v2/management-service/src/cmd/management/

ENTRYPOINT ["go","run","."]

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=build /detecctor-v2/management-service/management-service /usr/bin/detecctor-v2/management-service
WORKDIR /usr/bin/detecctor-v2/

ENTRYPOINT ["./management-service"]