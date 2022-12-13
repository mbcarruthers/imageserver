FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o imageserver ./cmd

RUN chmod +x /app/imageserver

CMD ["/app/imageserver"]