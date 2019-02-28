FROM arm32v7/golang as builder
WORKDIR /go/src/go-scape-market
COPY scraper.go .
RUN go get golang.org/x/net/html && \ 
    CGO_ENABLED=0 GOOS=linux go install go-scape-market

FROM arm32v7/ubuntu:18.04
LABEL maintainer="cgianelle@gmail.com"
LABEL version="0.0.1"
RUN apk --no-cache add ca-certificates
RUN adduser -S -D -H -h /app scraper
USER scraper
COPY --from=builder /go/bin/go-scape-market /app/
WORKDIR /app
CMD ["./go-scape-market"]

 
