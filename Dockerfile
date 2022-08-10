FROM golang:latest as builder

RUN mkdir -p /analytics
ADD . /analytics
WORKDIR /analytics

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /analytics ./cmd/main.go

FROM scratch
COPY --from=builder /analytics /analytics
COPY --from=builder /etc/ssl/certs /etc/ssl/certs/
WORKDIR /analytics

CMD ["./main"]
EXPOSE 3000 4000
