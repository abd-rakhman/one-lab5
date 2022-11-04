# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN apk add --no-cache bash
RUN apk add --no-cache libc6-compat

COPY main.go ./

RUN go build -o /register-crud

EXPOSE 3000

CMD [ "/register-crud" ]
