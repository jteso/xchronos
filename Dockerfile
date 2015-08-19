FROM golang:latest
MAINTAINER Javier Teso jtejob@gmail.com

RUN go get github.com/constabulary/gb/...

RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app
RUN gb build

CMD ["go-wrapper", "run"]
