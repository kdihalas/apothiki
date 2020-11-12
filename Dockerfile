FROM golang:1.14 as builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go build -o docker-repo cmd/*

FROM gcr.io/distroless/base
COPY --from=builder /go/src/app/docker-repo /
CMD ["/docker-repo"]