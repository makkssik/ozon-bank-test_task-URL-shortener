FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url-shortener ./cmd/url-shortener/main.go

FROM golang:1.26
WORKDIR /app

COPY --from=builder /app/url-shortener .
COPY --from=builder /app/config ./config
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./url-shortener", "-storage=postgres"]