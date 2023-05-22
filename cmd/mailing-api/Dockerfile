FROM golang:1.20.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY cmd/mailing-api/main.go cmd/mailing-api/main.go
COPY pkg/mailing pkg/mailing
COPY shared/ shared/
COPY cmd/mailing-api/.air.toml cmd/mailing-api/.air.toml


RUN go install github.com/cosmtrek/air@latest
RUN go build -o /app/main ./cmd/mailing-api/main.go


EXPOSE 8080

CMD ["air", "-c", "cmd/mailing-api/.air.toml"]

