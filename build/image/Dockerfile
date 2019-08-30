FROM michaelgrifalconi/git2go:alpine-golang-v1

WORKDIR /go/src/github.com/michaelgrifalconi/trufflegopher
COPY . / ./

ENV GOBIN $GOPATH/bin

RUN go test ./...

RUN go build -o tg ./cmd/tg/main.go
RUN mv tg /go/bin/tg

ENTRYPOINT ["tg"]
