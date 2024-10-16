# build-env
FROM golang:1.23.2-alpine AS build-env

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o main .

# final stage
FROM alpine:latest

ENV TZ=Asia/Taipei

RUN apk update && apk add -U tzdata

WORKDIR /app

COPY --from=build-env /app/main .

CMD ["/app/main"]