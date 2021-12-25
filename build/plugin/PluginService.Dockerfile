FROM golang:1.17-alpine as base

WORKDIR /detecctor-v2/plugin-service/src
COPY .. .

RUN mkdir "/detecctor-v2/plugins"
ARG PLUGIN_DIR
ENV PLUGIN_DIR=$PLUGIN_DIR

COPY $PLUGIN_DIR ../plugins

RUN go mod download
RUN go mod verify

RUN go build -o ../cmd-service ./cmd/cmd/main.go

FROM base as dev
WORKDIR /detecctor-v2/plugin-service/src/cmd/plugin
ENTRYPOINT ["go","run","."]

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=build /detecctor-v2/plugin-service/plugin-service /usr/bin/detecctor-v2/plugin-service
WORKDIR /usr/bin/detecctor-v2

ENTRYPOINT ["./plugin-service"]