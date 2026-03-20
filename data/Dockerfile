FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY data/go.mod data/go.sum ./data/

COPY pkg ./pkg

WORKDIR /app/data

RUN go mod download

COPY data/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/data/main .

EXPOSE 8103

CMD ["./main"]
