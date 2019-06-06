FROM golang:1.12 AS builder
WORKDIR /go/src/github.com/arbourd/audition/
COPY . ./

ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN go build -o app .

FROM node:8.2 AS webpack
WORKDIR /client
COPY client .
RUN yarn
RUN yarn build

FROM alpine:latest
WORKDIR /
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/arbourd/audition/app .
COPY --from=webpack client/dist ./client/dist
CMD ["./app"]
