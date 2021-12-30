FROM golang:1.17.5-alpine3.15 AS builder

WORKDIR /app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ADD . .

RUN go mod tidy
RUN go build -o /super-tiny-proxy

FROM scratch

COPY --from=builder /super-tiny-proxy /
