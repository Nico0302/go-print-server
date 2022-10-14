# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/cloudprint/main.go ./

RUN go build -o /cloudprint

EXPOSE 8080

CMD [ "/cloudprint" ]