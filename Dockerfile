FROM golang:1.11-alpine3.8
LABEL MAINTAINER="Nick Pleatsikas <nick@pleasikas.me>"

# For some reason this isn't set by default.
ENV GOPATH=/go

# Repository URL.
ENV REPO_URL=github.com/MrFlynn/thesaurize-bot

# Install git so remote libraries can be installed.
RUN apk update && \
    apk add git 

# Download and install app from remote.
RUN go get ${REPO_URL} && \
    go install ${REPO_URL}

# Run the app.
CMD thesaurize-bot