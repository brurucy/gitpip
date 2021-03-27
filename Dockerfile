FROM golang:alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN apk add --no-cache bash

EXPOSE 8080

CMD ["go", "run", "main.go"]
