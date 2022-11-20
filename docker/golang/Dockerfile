FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY ../.. .
RUN go mod download
RUN go build -o ./medusa ./src/main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/medusa .
EXPOSE 8000
ENTRYPOINT ["./medusa"]