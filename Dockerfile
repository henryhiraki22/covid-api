FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN pwd
RUN go build -o /build/main core/main.go

WORKDIR /dist

RUN cp /build/main .

EXPOSE 8030

CMD ["/dist/main"]