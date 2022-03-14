# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /factorial

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /factorial_serv