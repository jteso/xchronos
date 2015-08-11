FROM golang:latest
MAINTAINER Javier Teso jtejob@gmail.com

ENV package=github.com/jteso/xchronos
ENV executable=xchronos

RUN mkdir -p /go/src/$package
ADD . /go/src/$package
WORKDIR /go/src/$package
RUN go build -o $executable .
RUN chmod +x $executable

CMD ./$executable
