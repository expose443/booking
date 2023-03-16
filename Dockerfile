# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o cmd/hotel cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app .
LABEL version="1.0" 
LABEL creators="@abdu0222"
EXPOSE 9090
CMD [ "cmd/hotel" ]

