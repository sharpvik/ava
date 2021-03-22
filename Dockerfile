# Build server.
FROM golang:1.16-alpine3.12 AS server_builder
RUN mkdir /app
WORKDIR /app
# Add trusted certificates.
RUN apk --no-cache add ca-certificates
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on
RUN go build -o serve
# => /srv/serve*



# Migrate the build artifact, files, and dirs to Ubuntu.
FROM ubuntu:latest
RUN mkdir /ava
RUN mkdir /ava/storage
WORKDIR /ava
# Copying things over ...
COPY --from=server_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server_builder /app/serve /ava/serve
# At this point we will have:
#
#     /ava
#     --> storage/
#     --> serve*
#
# Exposing default HTTP port.
EXPOSE 80


# Start the server.
CMD [ "/ava/serve" ]