FROM golang:1.20-alpine

WORKDIR /usr/src/app
COPY . .
RUN go build -v -o /usr/local/bin/notification-service ./cmd/notification-service

CMD ["notification-service"]

