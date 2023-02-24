FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -o go-server
WORKDIR /dist
RUN cp /build/go-server .

FROM alpine
COPY --from=builder /dist/go-server .
EXPOSE 8888
ENTRYPOINT ["/go-server"]
# docker build . -t go-server