#builder
FROM golang:1.22-alpine as builder
RUN apk add --update make git protoc protobuf protobuf-dev curl

COPY Makefile /home/Makefile
COPY go.mod /home/go.mod
COPY go.sum /home/go.sum

WORKDIR /home
RUN make deps-go

COPY . /home
RUN make build-go

# gRPC Server

FROM alpine:latest as server

ARG GITHUB_PATH=github.com/arslanovdi/logistic-package-api

LABEL org.opencontainers.image.source=https://${GITHUB_PATH}

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /home/bin/grpc-server .
COPY --from=builder /home/config.yml .
COPY --from=builder /home/migrations/ ./migrations
COPY --from=builder /home/swagger ./swagger
COPY --from=builder /home/swagger-ui ./swagger-ui

RUN chown root:root grpc-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 9100

CMD ["./grpc-server"]
