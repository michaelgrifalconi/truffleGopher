FROM golang:alpine
#TODO: use fixed version -> immutable builds

RUN apk update && apk add --no-cache bash git openssh alpine-sdk cmake openssl-dev python

WORKDIR /go/src/github.com/michaelgrifalconi/trufflegopher
COPY . / ./

RUN ./build/image/setup-libgit2

ENV GOBIN $GOPATH/bin
#TODO: run all go tests before building
RUN go build -o tg ./cmd/tg/main.go
RUN mv tg /go/bin/tg

ENTRYPOINT ["tg"]
