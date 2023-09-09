default: run

build:
		go build -a -v -o pidgin-relay main.go

run:
		./pidgin-relay -c config.yaml
