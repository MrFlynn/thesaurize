FROM alpine:latest as base

# Create a user to copy over to target image.
RUN adduser -u 10000 -H -D porty

# Target container.
FROM scratch

# Want SSL certificates and users.
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /etc/passwd /etc/passwd

COPY bot /bin/bot

USER bot

ENTRYPOINT [ "/bin/bot" ]
CMD [ "--help" ]