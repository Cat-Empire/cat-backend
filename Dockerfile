FROM golang:1.17-alpine

RUN apk update && apk add make \ 
    rm -rf /var/cache/apk/*

RUN mkdir -p /usr/src/server
WORKDIR /usr/src/server

COPY . .

CMD ["make", "run"]
