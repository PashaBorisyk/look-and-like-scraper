FROM golang:1.12 as builder

LABEL maintainer="pashaborisyk <pashaborisyk@gmail.com>"

WORKDIR /go/src/look-and-like-scraper

COPY . .

RUN go get -d -v ./...

WORKDIR /go/src/look-and-like-scraper/main

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/look-and-like-scraper .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
    
COPY --from=builder /go/bin/look-and-like-scraper .

COPY ./resources ./resources

CMD ["./look-and-like-scraper"]