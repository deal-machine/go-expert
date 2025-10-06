FROM golang:alpine

WORKDIR /app

COPY . /app

RUN go mod tidy

RUN go build -o go-docker-cli

ENTRYPOINT ["./go-docker-cli", "stress"]