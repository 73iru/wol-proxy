FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /wol-proxy

EXPOSE 8089

CMD ["./wol-proxy"]