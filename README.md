# GoSniff

Meant to listen messages in discord channels and forward them to another.

# How to run locally

All you need is docker.

## Build service

```shell
make build
```

or if you doesn't have a make

```shell
docker build --tag discord--relay:latest .
```

## Run service

```shell
make run
```

or if you doesn't have a make

```shell
docker run -ti --rm -v `pwd`/config.yaml:/discord--relay/config.yaml discord--relay:latest
```

# Configuration

Fix your config.yaml

# General Flow
A single listener can listen to X number of channels across multiple servers. Each listener
opens up a connection to discord and then has to filter the messages only to what we care about.
Every new message fires a "message create" event that we must process. If we care about it, we clean
it up, transform the names and make it look nice and then send it off to the forwarder to relay it.
Forwarders have two modes currently, chatting as a user, and chatting as a bot using webhooks. Chatting
as a user requires less setup, but grants us far less control over the message. Using webhooks allow us 
to imitate it coming from the user that posted it. 
