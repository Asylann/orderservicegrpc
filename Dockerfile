FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/server

RUN go build -o order_service main.go

EXPOSE 50052

CMD ["./order_service"]