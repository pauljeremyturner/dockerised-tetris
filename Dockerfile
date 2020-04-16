FROM golang:1.13-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build ./...

CMD /app/run.sh