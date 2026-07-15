FROM golang:1.26-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/url-shortener .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/url-shortener /app/url-shortener
EXPOSE 8080
ENV HTTP_ADDR=:8080
CMD ["/app/url-shortener"]