# Gunakan image golang versi 1.21.1 sebagai base image
FROM golang:1.21.1

# Membuat direktori 'app' di dalam container
RUN mkdir /app

WORKDIR /app

COPY ./ .

RUN go get -d -v ./...

RUN go build -o ato-chat-app .

CMD ["./ato-chat-app"]
