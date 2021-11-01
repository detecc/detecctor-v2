FROM golang:1.17 as base
WORKDIR /detecctor-v2/management-service/src
COPY . .

FROM base as dev
RUN cd cmd/management/
ENTRYPOINT ["go","run","."]

FROM base as build
RUN go build main.go -o ../../../management-service

FROM base as test
RUN go test -v

FROM alpine as app

COPY --from=build /detecctor-v2/management-service/management-service /detecctor-v2/management-service
RUN mv  /detecctor-v2/management-service /usr/bin/management-service

ENTRYPOINT ["management-service"]