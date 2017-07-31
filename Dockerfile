FROM golang:1.8.3 AS builder
WORKDIR /go/src/github.com/arbourd/audition/
COPY *.go ./
COPY vendor ./vendor
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM node:8.2 AS webpack
WORKDIR /client
COPY client .
RUN yarn
RUN yarn build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/arbourd/audition/app .
COPY --from=webpack client/dist ./client/dist
CMD ["./app"]
