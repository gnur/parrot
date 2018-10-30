#build javascript
FROM node as jsbuilder
WORKDIR /workspace
COPY app /workspace
RUN yarn
RUN yarn build


#pack static files and create binary
FROM golang:1.11-alpine as builder
WORKDIR /go/src/github.com/gnur/parrot/
RUN apk --no-cache add git

RUN go get github.com/GeertJohan/go.rice
RUN go get github.com/GeertJohan/go.rice/rice
COPY --from=jsbuilder /workspace/dist app/dist
COPY vendor vendor
COPY app app
COPY cmd cmd
COPY pkg pkg
RUN rice embed-go -v -i ./pkg/webserver
RUN go build -ldflags="-s -w" -o parrot ./cmd/parrot
RUN go build -ldflags="-s -w" -o log-generator ./cmd/log-generator


#actual container
FROM alpine
COPY --from=builder /go/src/github.com/gnur/parrot/parrot /parrot
COPY --from=builder /go/src/github.com/gnur/parrot/log-generator /log-generator

ENTRYPOINT [ "/parrot" ]
CMD [ "-listen", "udp://0.0.0.0:514", "-web", "0.0.0.0:8080" ]
