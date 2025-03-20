# Dockerfile
FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN go build -o testapi .

CMD ["./testapi"]