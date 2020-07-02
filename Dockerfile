FROM golang:1.14 as build

# Set version, commit, and date information at compile time.
ARG VERSION
ARG COMMIT
ARG DATE

# Standard environment variables.
ENV CGO_ENABLED=0
ENV GOOS=linux

# Create a user to copy over to target image.
RUN useradd -u 10000 bot

WORKDIR /go/src
COPY . .

# Build the executable.
RUN go build -ldflags="-s -w -X 'main.Version=${VERSION}' -X 'main.Version=${COMMIT}' -X 'main.date=${DATE}'" -o bot ./cmd/bot

# Target container.
FROM scratch

# Want SSL certificates and users.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /go/src/bot /bin/bot

USER bot

ENTRYPOINT [ "/bin/bot" ]
CMD [ "--help" ]