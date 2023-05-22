FROM golang:1.20.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY cmd/newsletter-api/main.go cmd/newsletter-api/main.go
COPY pkg/newsletter pkg/newsletter
COPY shared/ shared/
COPY cmd/newsletter-api/.air.toml cmd/newsletter-api/.air.toml

RUN go install github.com/cosmtrek/air@latest
RUN go build -o /app/main ./cmd/newsletter-api/main.go


EXPOSE 8080

CMD ["air", "-c", "cmd/newsletter-api/.air.toml"]

