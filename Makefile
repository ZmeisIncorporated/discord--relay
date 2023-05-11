default: run

run:
	docker run -ti --rm -v `pwd`/config.yaml:/discord--relay/config.yaml discord--relay:latest


build:
	docker build --tag discord--relay:latest .
