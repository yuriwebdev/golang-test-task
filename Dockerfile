FROM ubuntu
MAINTAINER yurigasparyan <yurwebdev.yur@gmail.com>
RUN apt-get update
RUN apt-get install nano
RUN apt-get install -y software-properties-common python-software-properties
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get -y install golang-go git
RUN mkdir /work
ENV GOPATH=/work

RUN go get github.com/yuriwebdev/golang-test-task
RUN go get github.com/PuerkitoBio/goquery
RUN go build github.com/yuriwebdev/golang-test-task

CMD /golang-test-task -addr $ADDR

