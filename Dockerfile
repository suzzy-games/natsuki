# Build Application
FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Copy Built Application
FROM alpine:3.16.0 AS runtime
WORKDIR /app
COPY --from=builder /app/main .

# Run Application
EXPOSE 80
EXPOSE 443
CMD [ "/app/main" ]