FROM golang:1.12-alpine AS builder
WORKDIR /go/src/github.com/arbourd/audition/
RUN apk --no-cache add --update git

ENV CGO_ENABLED 0
ENV GO111MODULE on

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o app .

FROM node:12-alpine AS webpack
WORKDIR /client

COPY client/package.json client/yarn.lock ./
RUN yarn

COPY client .
RUN yarn build

FROM alpine:latest
WORKDIR /
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/arbourd/audition/app .
COPY --from=webpack client/dist ./client/dist
CMD ["./app"]
