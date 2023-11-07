default: build

build:
		go build -a -v -o pidgin-relay main.go

install:
		rm -f /usr/local/bin/pidgin-relay
		cp -v pidgin-relay /usr/local/bin/

run:
		./pidgin-relay -c config.yaml

test:
		go run main.go -c config.yaml
