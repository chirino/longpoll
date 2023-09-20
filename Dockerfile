FROM golang:1.20-alpine as builder

WORKDIR /src
COPY go.mod .
#COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server ./server
RUN CGO_ENABLED=0 go build -o /go/bin/client ./client

FROM alpine:3.16
COPY --from=builder /go/bin/server /server
COPY --from=builder /go/bin/client /client

ENTRYPOINT [ "/server" ]
