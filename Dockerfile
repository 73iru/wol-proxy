FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go /app/

RUN CGO_ENABLED=0  GOOS=linux go build main.go

EXPOSE 8089
EXPOSE 9

CMD ["./main"]