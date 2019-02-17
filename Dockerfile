FROM golang as builder
WORKDIR /go/src/go-scape-market
COPY scraper.go .
RUN go get golang.org/x/net/html && \ 
    CGO_ENABLED=0 GOOS=linux go install go-scape-market

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -S -D -H -h /app scraper
USER scraper
COPY --from=builder /go/bin/go-scape-market /app/
WORKDIR /app
CMD ["./go-scape-market"]

 
