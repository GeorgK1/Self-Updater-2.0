# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./

ENV GIT_USERNAME = "your name here"
ENV GIT_PASSWORD = "your git token here"
ENV GIT_URL = "your git url here"
ENV GIT_REMOTE = "main"

RUN go mod tidy
RUN go build

EXPOSE 8080

