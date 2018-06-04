FROM golang:1.10.2-alpine3.7 as builder
COPY go.api /go/src/app
WORKDIR /go/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.6
WORKDIR /
COPY --from=builder /go/src/app/go.api /
RUN adduser -s /bin/false -S -D -H app
USER app
EXPOSE 6080
CMD ["/go.api"]
