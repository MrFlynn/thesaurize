FROM golang:1.14 as build

# Create a user to copy over to target image.
RUN useradd -u 10000 bot

# Target container.
FROM scratch

# Want SSL certificates and users.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd

COPY bot /bin/bot

USER bot

ENTRYPOINT [ "/bin/bot" ]
CMD [ "--help" ]