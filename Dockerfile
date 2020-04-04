FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get -u github.com/cosmtrek/air

CMD air