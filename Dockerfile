FROM golang:1.11-alpine3.8
LABEL MAINTAINER="Nick Pleatsikas <nick@pleasikas.me>"

# For some reason this isn't set by default.
ENV GOPATH=/go

# Install git so remote libraries can be installed.
RUN apk update && \
    apk add git 

COPY . /go/src/github.com/MrFlynn/thesaurus-bot
WORKDIR /go/src/github.com/MrFlynn/thesaurus-bot

# Install dependency.
RUN go get github.com/bwmarrin/discordgo

# Build and run.
RUN go build
CMD ./thesaurus-bot