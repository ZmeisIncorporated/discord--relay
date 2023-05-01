FROM golang:1.19

WORKDIR /discord--relay
COPY . .

RUN go build -a -v -o app main.go

ENTRYPOINT ["/discord--relay/app"]
