FROM golang:1.24.5 AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

COPY --from=builder /app/app . 
EXPOSE 8080

ENTRYPOINT [ "./app" ]

