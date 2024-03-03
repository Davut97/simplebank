#Binary build
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Final image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY app.env .


EXPOSE 8080
CMD ["/app/main"]
