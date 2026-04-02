#build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service main.go
RUN adduser -D -u 10001 appuser

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/service /service
USER appuser
ENV PORT=8000
ENV RESPONSE_MESSAGE="Service request succeeded!"
ENV ALLOW_ORIGIN="*"
EXPOSE 8000
ENTRYPOINT ["/service"]