FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["/app/main"]