FROM golang:1.17 as base
WORKDIR /detecctor-v2/plugin-service/src
COPY . .
RUN mkdir "/detecctor-v2/plugins"
WORKDIR /detecctor-v2/plugin-service/src/cmd/plugin


FROM base as dev
ENTRYPOINT ["go","run","."]

FROM base as build

ARG PLUGIN_DIR
ENV PLUGIN_DIR=$PLUGIN_DIR

COPY $PLUGIN_DIR ../plugins

RUN go build main.go -o ../../../plugin-service

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=build /detecctor-v2/plugin-service/plugin-service /detecctor-v2/plugin-service
RUN mv /detecctor-v2/plugin-service /usr/bin/plugin-service
ENTRYPOINT ["plugin-service"]