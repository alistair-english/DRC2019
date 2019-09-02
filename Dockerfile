FROM gocv

RUN apt-get update
RUN apt-get install nano

WORKDIR /go/src/github.com/alistair-english/DRC2019

COPY . /go/src/github.com/alistair-english/DRC2019

RUN go get /go/src/github.com/alistair-english/DRC2019/cmd/main

RUN rm -rf /go/src/github.com/alistair-english/DRC2019