#build stage
FROM golang:1.21-alpine AS builder
RUN adduser -D -u 10001 nonroot
WORKDIR /devops-service
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /etc/passwd /etc/passwd
# path fixed: Match the WORKDIR from builder
COPY --from=builder /devops-service/service /service

USER nonroot

ENV PORT=8000
ENV RESPONSE_MESSAGE="Service request succeeded!"
ENV ALLOW_ORIGIN="*"

EXPOSE 8000

ENTRYPOINT ["/service"]