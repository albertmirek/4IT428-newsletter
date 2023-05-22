FROM golang:1.20.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY cmd/user-api/main.go cmd/user-api/main.go
COPY pkg/user pkg/user
COPY shared/ shared/
COPY cmd/user-api/.air.toml cmd/user-api/.air.toml

RUN go install github.com/cosmtrek/air@latest
RUN go build -o /app/main ./cmd/user-api/main.go


EXPOSE 8080

CMD ["air", "-c", "cmd/user-api/.air.toml"]

# CMD if [ "$ENVIRONMENT" = "production" ] ; then ./user-api ; else air ; fi

# CMD ["./cmd/user-api/main.go"]


