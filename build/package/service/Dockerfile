FROM golang:1.19.3-alpine AS builder
WORKDIR /build
ARG service
COPY . .
WORKDIR /build/cmd/$service
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.17
WORKDIR /root/
ARG service
COPY --from=builder /build/cmd/$service/$service service
CMD ["./service"]
