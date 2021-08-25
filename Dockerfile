FROM golang:1.16

WORKDIR /workdir
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app cmd/*.go

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /bin
COPY --from=0 /workdir/bin/app .
EXPOSE 8080

ENTRYPOINT ["app"]