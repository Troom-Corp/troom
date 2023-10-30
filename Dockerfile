FROM golang:latest

COPY . /app
WORKDIR /app/cmd/troom

RUN go build -o ./server

RUN go mod tidy
RUN go mod download

EXPOSE 5000

CMD ["./server"]