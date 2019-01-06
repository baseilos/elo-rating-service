FROM golang:alpine AS builder

COPY . $GOPATH/src/jozeflang.com/elo-rating-service/
WORKDIR $GOPATH/src/jozeflang.com/elo-rating-service/

RUN go build -o /go/bin/elo-rating-service

FROM alpine
COPY --from=builder /go/bin/elo-rating-service /go/bin/elo-rating-service
ENTRYPOINT /go/bin/elo-rating-service