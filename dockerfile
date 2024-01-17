FROM golang:1.21.1

# membuat direktori app
RUN mkdir /app

# set working directory /app
WORKDIR /app

COPY ./ /app

RUN go get ./...

RUN go build -o ato-chat-app

CMD ["ato-chat-app"]